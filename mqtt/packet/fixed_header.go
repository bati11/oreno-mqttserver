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
	result := FixedHeader{packetType, dup, qos1, qos2, retain, 0}
	return result
}

func refbit(i byte, b uint) byte {
	return (i >> b) & 1
}
