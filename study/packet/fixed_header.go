package packet

type FixedHeader struct {
	PacketType byte
}

func ToFixedHeader(bs []byte) FixedHeader {
	b := bs[0]
	packetType := b >> 4
	return FixedHeader{packetType}
}
