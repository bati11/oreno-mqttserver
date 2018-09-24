package packet

import "fmt"

type ConnectFlags struct {
}

type ConnectVariableHeader struct {
	ProtocolName  string
	ProtocolLevel uint8
	ConnectFlags  byte
	KeepAlive     uint16
}

func ToConnectVariableHeader(fixedHeader FixedHeader, bs []byte) (ConnectVariableHeader, error) {
	if fixedHeader.PacketType() != CONNECT {
		return ConnectVariableHeader{}, fmt.Errorf("packet type is invalid. it got is %v", fixedHeader.PacketType())
	}

	protocolName := bs[:6]
	if !isValidProtocolName(protocolName) {
		return ConnectVariableHeader{}, fmt.Errorf("protocol name is invalid. it got is %q", protocolName)
	}

	protocolLevel := bs[6]
	if protocolLevel != 4 {
		return ConnectVariableHeader{}, fmt.Errorf("protocol level is not supported. it got is %v", protocolLevel)
	}

	connectFlags := bs[7]
	reserved := connectFlags & 1
	if reserved != 0 {
		return ConnectVariableHeader{}, fmt.Errorf("reserved value in connect flags must be 0. it got is %v", reserved)
	}

	result := ConnectVariableHeader{
		ProtocolName:  "MQTT",
		ProtocolLevel: 4,
	}
	return result, nil
}

func isValidProtocolName(protocolName []byte) bool {
	return len(protocolName) == 6 &&
		protocolName[0] == 0 && protocolName[1] == 4 &&
		protocolName[2] == 'M' && protocolName[3] == 'Q' && protocolName[4] == 'T' && protocolName[5] == 'T'
}
