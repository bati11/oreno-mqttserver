/*
Package packet implements packets of MQTT.

*/
package packet

import (
	"encoding/binary"
	"errors"
)

type PacketType uint8

const (
	Reserved PacketType = iota
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
	Reserved2
)

func (p PacketType) string() string {
	switch p {
	case Reserved:
		return "Reserved"
	case CONNECT:
		return "CONNECT"
	default:
		return "unknown"
	}
}

func (p PacketType) byte() byte {
	return byte(p)
}

// FixedHeader is part of MQTT Control Packet.
type FixedHeader struct {
	PacketType      PacketType
	Dup             bool
	QoS1            bool
	QoS2            bool
	Retain          bool
	RemainingLength uint
}

var (
	ErrBytesLength     = errors.New("fixed header bytes length should be => 2")
	ErrPacketTypeValue = errors.New("packet type is between 0 and 15")
)

// ToFixedHeader converts bytes into a FixedHeader structure.
func ToFixedHeader(bs []byte) (FixedHeader, []byte, error) {
	if len(bs) < 2 {
		return FixedHeader{}, nil, ErrBytesLength
	}
	packetType := bs[0] >> 4
	if packetType < 0 || 15 < packetType {
		return FixedHeader{}, nil, ErrBytesLength
	}
	dup := refbit(bs[0], 3) > 0
	qos1 := refbit(bs[0], 2) > 0
	qos2 := refbit(bs[0], 1) > 0
	retain := refbit(bs[0], 0) > 0
	remainingLength, remains := decodeRemainingLength(bs[1:])
	result := FixedHeader{
		PacketType:      PacketType(packetType),
		Dup:             dup,
		QoS1:            qos1,
		QoS2:            qos2,
		Retain:          retain,
		RemainingLength: remainingLength,
	}
	return result, remains, nil
}

func (h *FixedHeader) ToBytes() []byte {
	var result []byte
	b := h.PacketType.byte() << 4
	if h.Dup {
		b = onbit(b, 4)
	}
	if h.QoS1 {
		b = onbit(b, 2)
	}
	if h.QoS2 {
		b = onbit(b, 1)
	}
	if h.Retain {
		b = onbit(b, 0)
	}
	result = append(result, b)
	remainingLength := h.encodeRemainingLength()
	result = append(result, remainingLength...)
	return result
}

func refbit(b byte, i uint) byte {
	return (b >> i) & 1
}

func onbit(b byte, i uint) byte {
	return b | (1 << i)
}

func decodeRemainingLength(bs []byte) (uint, []byte) {
	multiplier := uint(1)
	var value uint
	i := uint(0)
	for ; i < 8; i++ {
		digit := bs[i]
		value = value + uint(digit&127)*multiplier
		multiplier = multiplier * 128
		if (digit & 128) == 0 {
			break
		}
	}
	remains := bs[i+1:]
	return value, remains
}

func (h *FixedHeader) encodeRemainingLength() []byte {
	x := h.RemainingLength
	var encodedByte uint64
	for {
		encodedByte = uint64(x % 128)
		x = x / 128
		if x > 0 {
			encodedByte = encodedByte | 128
		} else {
			break
		}
	}
	buf := make([]byte, binary.MaxVarintLen32)
	n := binary.PutUvarint(buf, encodedByte)
	return buf[:n]
}
