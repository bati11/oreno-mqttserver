package packet

import (
	"encoding/binary"
	"fmt"
)

type ConnectFlags struct {
	CleanSession bool
	WillFlag     bool
	WillQoS      byte
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

func ToConnectVariableHeader(fixedHeader FixedHeader, bs []byte) (ConnectVariableHeader, []byte, error) {
	if fixedHeader.PacketType != CONNECT {
		return ConnectVariableHeader{}, nil, fmt.Errorf("packet type is invalid. it got is %v", fixedHeader.PacketType)
	}

	protocolName := bs[:6]
	if !isValidProtocolName(protocolName) {
		return ConnectVariableHeader{}, nil, fmt.Errorf("protocol name is invalid. it got is %q", protocolName)
	}

	protocolLevel := bs[6]
	if protocolLevel != 4 {
		return ConnectVariableHeader{}, nil, fmt.Errorf("protocol level is not supported. it got is %v", protocolLevel)
	}

	connectFlagsBytes := bs[7]
	reserved := connectFlagsBytes & 1
	if reserved != 0 {
		return ConnectVariableHeader{}, nil, fmt.Errorf("reserved value in connect flags must be 0. it got is %v", reserved)
	}

	connectFlags := ConnectFlags{}

	// TODO now, support only 1
	cleanSession := refbit(connectFlagsBytes, 1)
	if cleanSession != 1 {
		return ConnectVariableHeader{}, nil, fmt.Errorf("clean session value in connect flags must be 1. it got is %v", cleanSession)
	}
	connectFlags.CleanSession = true

	// TODO now, support only 0
	willFlag := refbit(connectFlagsBytes, 2)
	if willFlag != 0 {
		return ConnectVariableHeader{}, nil, fmt.Errorf("will flag value in connect flags must be 0. it got is %v", willFlag)
	}
	connectFlags.WillFlag = false

	// TODO now, support only QoS0
	willQoS2 := refbit(connectFlagsBytes, 3)
	willQoS1 := refbit(connectFlagsBytes, 4)
	if willQoS2 != 0 || willQoS1 != 0 {
		return ConnectVariableHeader{}, nil, fmt.Errorf("will QoS value in connect flags must be 0. it got is %v %v", willQoS1, willQoS2)
	}
	connectFlags.WillQoS = 0

	// TODO now, support only 0
	willRetain := refbit(connectFlagsBytes, 5)
	if willRetain != 0 {
		return ConnectVariableHeader{}, nil, fmt.Errorf("will retain value in connect flags must be 0. it got is %v", willRetain)
	}
	connectFlags.WillRetain = false

	// TODO now, support only 0
	passwordFlag := refbit(connectFlagsBytes, 6)
	if passwordFlag != 0 {
		return ConnectVariableHeader{}, nil, fmt.Errorf("password flag value in connect flags must be 0. it got is %v", passwordFlag)
	}
	connectFlags.PasswordFlag = false

	// TODO now, support only 0
	userNameFlag := refbit(connectFlagsBytes, 7)
	if userNameFlag != 0 {
		return ConnectVariableHeader{}, nil, fmt.Errorf("user name flag value in connect flags must be 0. it got is %v", userNameFlag)
	}
	connectFlags.UserNameFlag = false

	keepAlive := binary.BigEndian.Uint16(bs[8:10])

	result := ConnectVariableHeader{
		ProtocolName:  "MQTT",
		ProtocolLevel: 4,
		ConnectFlags:  connectFlags,
		KeepAlive:     keepAlive,
	}
	return result, bs[10:], nil
}

func isValidProtocolName(protocolName []byte) bool {
	return len(protocolName) == 6 &&
		protocolName[0] == 0 && protocolName[1] == 4 &&
		protocolName[2] == 'M' && protocolName[3] == 'Q' && protocolName[4] == 'T' && protocolName[5] == 'T'
}
