package packet

import (
	"encoding/binary"
)

type SubscribeVariableHeader struct {
	PacketIdentifier uint16
}

func (s *SubscribeVariableHeader) Length() uint {
	return 2 // PacketIdentifier (uint16) byte size
}

func (reader *MQTTReader) readSubscribeVariableHeader() (*SubscribeVariableHeader, error) {
	packetIdentifierMSB, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	packetIdentifierLSB, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	packetIdentifier := binary.BigEndian.Uint16([]byte{packetIdentifierMSB, packetIdentifierLSB})

	return &SubscribeVariableHeader{packetIdentifier}, nil
}
