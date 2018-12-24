package packet

type FixedHeader struct {
	PacketType      byte
	Dup             byte
	QoS1            byte
	QoS2            byte
	Retain          byte
	RemainingLength uint
}

func ToFixedHeader(bs []byte) FixedHeader {
	b := bs[0]
	packetType := b >> 4
	dup := refbit(bs[0], 3)
	qos1 := refbit(bs[0], 2)
	qos2 := refbit(bs[0], 1)
	retain := refbit(bs[0], 0)
	remainingLength := decodeRemainingLength(bs[1:])
	return FixedHeader{
		PacketType:      packetType,
		Dup:             dup,
		QoS1:            qos1,
		QoS2:            qos2,
		Retain:          retain,
		RemainingLength: remainingLength,
	}
}

func refbit(b byte, n uint) byte {
	return (b >> n) & 1
}

// a
func decodeRemainingLength(bs []byte) uint {
	multiplier := uint(1)
	var value uint
	i := uint(0)
	for ; i < 8; i++ {
		b := bs[i]
		digit := b
		value = value + uint(digit&127)*multiplier
		multiplier = multiplier * 128
		if (digit & 128) == 0 {
			break
		}
	}
	return value
}
