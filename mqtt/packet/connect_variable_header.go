package packet

import (
	"bufio"
	"io"

	"github.com/pkg/errors"
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

func ToConnectVariableHeader(fixedHeader DefaultFixedHeader, r *bufio.Reader) (ConnectVariableHeader, error) {
	if fixedHeader.PacketType != 1 {
		return ConnectVariableHeader{}, errors.New("fixedHeader.PacketType must be 1")
	}
	protocolName := make([]byte, 6)
	_, err := io.ReadFull(r, protocolName)
	if err != nil || !isValidProtocolName(protocolName) {
		return ConnectVariableHeader{}, RefusedByUnacceptableProtocolVersion("protocol name is invalid")
	}
	protocolLevel, err := r.ReadByte()
	if err != nil || protocolLevel != 4 {
		return ConnectVariableHeader{}, RefusedByUnacceptableProtocolVersion("protocol level must be 4")
	}

	// TODO
	_, err = r.ReadByte() // connectFlags
	if err != nil {
		return ConnectVariableHeader{}, err
	}
	_, err = r.ReadByte() // keepAlive MSB
	if err != nil {
		return ConnectVariableHeader{}, err
	}
	_, err = r.ReadByte() // keepAlive LSB
	if err != nil {
		return ConnectVariableHeader{}, err
	}

	return ConnectVariableHeader{
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
