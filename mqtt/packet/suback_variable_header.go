package packet

import "encoding/binary"

type SubackVariableHeader struct {
	PacketIdentifier uint16
}

func (s *SubackVariableHeader) Length() uint {
	return 2 // uint16 size
}

func (s *SubackVariableHeader) ToBytes() []byte {
	result := make([]byte, s.Length())
	binary.BigEndian.PutUint16(result, s.PacketIdentifier)
	return result
}
