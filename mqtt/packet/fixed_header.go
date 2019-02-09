package packet

import (
	"bufio"
)

const (
	_ = iota
	CONNECT
	CONNACK
	PUBLISH
	PUBACK
	PUBREC
	PUBREL
	PUBCOMP
	SUBSCRIBE
	SUBACK
	UNSUBSCRIBE
	UNSUBACK
	PINGREQ
	PINGRESP
	DISCONNECT
)

type FixedHeader struct {
	PacketType      byte
	Dup             byte
	QoS1            byte
	QoS2            byte
	Retain          byte
	RemainingLength uint
}

func (h FixedHeader) ToBytes() []byte {
	var result []byte
	b := h.PacketType << 4
	result = append(result, b)
	remainingLength := encodeRemainingLength(h.RemainingLength)
	result = append(result, remainingLength...)
	return result
}

func ToFixedHeader(r *bufio.Reader) (FixedHeader, error) {
	b, err := r.ReadByte()
	if err != nil {
		return FixedHeader{}, err
	}
	packetType := b >> 4
	dup := refbit(b, 3)
	qos1 := refbit(b, 2)
	qos2 := refbit(b, 1)
	retain := refbit(b, 0)
	remainingLength, err := decodeRemainingLength(r)
	if err != nil {
		return FixedHeader{}, err
	}
	return FixedHeader{
		PacketType:      packetType,
		Dup:             dup,
		QoS1:            qos1,
		QoS2:            qos2,
		Retain:          retain,
		RemainingLength: remainingLength,
	}, nil
}

func refbit(b byte, n uint) byte {
	return (b >> n) & 1
}

// a
func decodeRemainingLength(r *bufio.Reader) (uint, error) {
	multiplier := uint(1)
	var value uint
	i := uint(0)
	for ; i < 8; i++ {
		b, err := r.ReadByte()
		if err != nil {
			return 0, err
		}
		digit := b
		value = value + uint(digit&127)*multiplier
		multiplier = multiplier * 128
		if (digit & 128) == 0 {
			break
		}
	}
	return value, nil
}

func encodeRemainingLength(x uint) []byte {
	var encodedByte byte
	var result []byte
	for {
		encodedByte = byte(x % 128)
		x = x / 128
		if x > 0 {
			encodedByte = encodedByte | 128
		}
		result = append(result, encodedByte)
		if x <= 0 {
			break
		}
	}
	return result
}
