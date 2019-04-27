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
	Reserved        byte
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

func (reader *MQTTReader) readFixedHeader() (*FixedHeader, error) {
	packetType, err := reader.ReadPacketType()
	if err != nil {
		return nil, err
	}
	reserved := *reader.byte1 >> 4
	remainingLength, err := decodeRemainingLength(reader.r)
	if err != nil {
		return nil, err
	}
	return &FixedHeader{
		PacketType:      packetType,
		Reserved:        reserved,
		RemainingLength: remainingLength,
	}, nil
}

type PublishFixedHeader struct {
	PacketType      byte
	Dup             byte
	QoS1            byte
	QoS2            byte
	Retain          byte
	RemainingLength uint
}

func NewPublishFixedHeader(packetType byte, remainingLength uint) *PublishFixedHeader {
	return &PublishFixedHeader{
		PacketType:      packetType,
		RemainingLength: remainingLength,
	}
}

func (h PublishFixedHeader) ToBytes() []byte {
	var result []byte
	b := h.PacketType << 4
	result = append(result, b)
	remainingLength := encodeRemainingLength(h.RemainingLength)
	result = append(result, remainingLength...)
	return result
}

func (reader *MQTTReader) readPublishFixedHeader() (*PublishFixedHeader, error) {
	packetType, err := reader.ReadPacketType()
	if err != nil {
		return nil, err
	}
	dup := refbit(*reader.byte1, 3)
	qos1 := refbit(*reader.byte1, 2)
	qos2 := refbit(*reader.byte1, 1)
	retain := refbit(*reader.byte1, 0)
	remainingLength, err := decodeRemainingLength(reader.r)
	if err != nil {
		return nil, err
	}
	return &PublishFixedHeader{
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
