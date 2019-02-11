package packet

import (
	"bufio"
	"encoding/binary"
	"fmt"
)

type SubscribeVariableHeader struct {
	PacketIdentifier uint16
}

func (s *SubscribeVariableHeader) Length() uint {
	return 2 // PacketIdentifier (uint16) byte size
}

func ToSubscribeVariableHeader(fixedHeader FixedHeader, r *bufio.Reader) (SubscribeVariableHeader, error) {
	if fixedHeader.PacketType != SUBSCRIBE {
		return SubscribeVariableHeader{}, fmt.Errorf("packet type is invalid. it got is %v", fixedHeader.PacketType)
	}

	packetIdentifierMSB, err := r.ReadByte()
	if err != nil {
		return SubscribeVariableHeader{}, err
	}
	packetIdentifierLSB, err := r.ReadByte()
	if err != nil {
		return SubscribeVariableHeader{}, err
	}
	packetIdentifier := binary.BigEndian.Uint16([]byte{packetIdentifierMSB, packetIdentifierLSB})

	return SubscribeVariableHeader{packetIdentifier}, nil
}
