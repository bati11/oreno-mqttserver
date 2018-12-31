package packet

import (
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

func ToConnectVariableHeader(fixedHeader FixedHeader, bs []byte) (ConnectVariableHeader, error) {
	if fixedHeader.PacketType != 1 {
		return ConnectVariableHeader{}, errors.New("fixedHeader.PacketType must be 1")
	}
	if !isValidProtocolName(bs[:6]) {
		return ConnectVariableHeader{}, errors.New("protocol name is invalid")
	}
	if bs[6] != 4 {
		return ConnectVariableHeader{}, errors.New("protocol level must be 4")
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
