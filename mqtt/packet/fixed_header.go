/*
Package packet implements packets of MQTT.

*/
package packet

import "errors"

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

// FixedHeader is part of MQTT Control Packet.
type FixedHeader struct {
	byte1           byte
	remainingLength uint
}

func (h FixedHeader) PacketType() PacketType {
	b := h.byte1 >> 4
	return PacketType(b)
}

func (h FixedHeader) Dup() bool {
	return refbit(h.byte1, 3) > 0
}

func (h FixedHeader) QoS1() bool {
	return refbit(h.byte1, 2) > 0
}

func (h FixedHeader) QoS2() bool {
	return refbit(h.byte1, 1) > 0
}

func (h FixedHeader) Retain() bool {
	return refbit(h.byte1, 0) > 0
}

func (h FixedHeader) RemainingLength() uint {
	return h.remainingLength
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
	remainingLength, remains := decodeRemainingLength(bs[1:])
	result := FixedHeader{bs[0], remainingLength}
	return result, remains, nil
}

func refbit(i byte, b uint) byte {
	return (i >> b) & 1
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
