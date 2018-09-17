package packet

type FixedHeader struct {
	PacketType      uint8
	Dup             byte
	QoS1            byte
	QoS2            byte
	Retain          byte
	RemainingLength uint8
}

func ToFixedHeader(bs [2]byte) FixedHeader {
	packetType := bs[0] >> 4
	dup := refbit(bs[0], 3)
	qos1 := refbit(bs[0], 2)
	qos2 := refbit(bs[0], 1)
	retain := refbit(bs[0], 0)
	result := FixedHeader{packetType, dup, qos1, qos2, retain, decodeRemainingLength(bs[1])}
	return result
}

func refbit(i byte, b uint) byte {
	return (i >> b) & 1
}

func decodeRemainingLength(b byte) uint8 {
	multiplier := uint8(1)
	var value uint8
	for i := uint(0); i < 8; i++ {
		digit := b
		value = value + uint8(digit&127)*multiplier
		multiplier = multiplier * 128
		if (digit & 128) == 0 {
			break
		}
	}
	return value
}
