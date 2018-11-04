package packet

import (
	"encoding/binary"
	"fmt"
)

type SubscribeVariableHeader struct {
	PacketIdentifier uint16
}

func ToSubscribeVariableHeader(bs []byte) (SubscribeVariableHeader, error) {
	if len(bs) < 2 {
		return SubscribeVariableHeader{}, fmt.Errorf("len(bs) should be >= 2")
	}
	packetIdentifier := binary.BigEndian.Uint16(bs[:2])
	result := SubscribeVariableHeader{packetIdentifier}
	return result, nil
}
