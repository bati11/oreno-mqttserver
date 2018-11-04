package packet

import "encoding/binary"

type SubackVariableHeader struct {
	PacketIdentifier uint16
}

func (h *SubackVariableHeader) ToBytes() []byte {
	result := make([]byte, 2)
	binary.BigEndian.PutUint16(result, h.PacketIdentifier)
	return result
}
