package packet

import "fmt"

type ConnectVariableHeader struct {
	protocolName  [6]byte
	protocolLevel byte
	connectFlags  byte
	keepAlive     [2]byte
}

func (h ConnectVariableHeader) ProtocolName() string {
	return "MQTT"
}

func (h ConnectVariableHeader) ProtocolLevel() uint8 {
	return 4
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

	result := ConnectVariableHeader{}
	copy(result.protocolName[:], protocolName)
	result.protocolLevel = protocolLevel
	return result, nil
}

func isValidProtocolName(protocolName []byte) bool {
	return len(protocolName) == 6 &&
		protocolName[0] == 0 && protocolName[1] == 4 &&
		protocolName[2] == 'M' && protocolName[3] == 'Q' && protocolName[4] == 'T' && protocolName[5] == 'T'
}
