package packet

import (
	"io"
)

type ConnectFlags struct {
	CleanSession bool
	WillFlag     bool
	WillQoS      uint8
	WillRetain   bool
	PasswordFlag bool
	UserNameFlag bool
}

type ConnectVariableHeader struct {
	ProtocolName  string
	ProtocolLevel uint8
	ConnectFlags  ConnectFlags
	KeepAlive     uint16
}

func (reader *MQTTReader) readConnectVariableHeader() (*ConnectVariableHeader, error) {
	protocolName := make([]byte, 6)
	_, err := io.ReadFull(reader.r, protocolName)
	if err != nil || !isValidProtocolName(protocolName) {
		return nil, RefusedByUnacceptableProtocolVersion("protocol name is invalid")
	}
	protocolLevel, err := reader.r.ReadByte()
	if err != nil || protocolLevel != 4 {
		return nil, RefusedByUnacceptableProtocolVersion("protocol level must be 4")
	}

	// TODO
	_, err = reader.r.ReadByte() // connectFlags
	if err != nil {
		return nil, err
	}
	_, err = reader.r.ReadByte() // keepAlive MSB
	if err != nil {
		return nil, err
	}
	_, err = reader.r.ReadByte() // keepAlive LSB
	if err != nil {
		return nil, err
	}

	return &ConnectVariableHeader{
		ProtocolName:  "MQTT",
		ProtocolLevel: 4,
		ConnectFlags:  ConnectFlags{UserNameFlag: true, PasswordFlag: true, WillRetain: false, WillQoS: 1, WillFlag: true, CleanSession: true},
		KeepAlive:     10,
	}, nil
}

func isValidProtocolName(protocolName []byte) bool {
	return len(protocolName) == 6 &&
		protocolName[0] == 0 && protocolName[1] == 4 &&
		protocolName[2] == 'M' && protocolName[3] == 'Q' && protocolName[4] == 'T' && protocolName[5] == 'T'
}
