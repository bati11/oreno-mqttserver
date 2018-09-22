/*
Package packet implements packets of MQTT.

*/
package packet

import "errors"

// FixedHeader is part of MQTT Control Packet.
type FixedHeader struct {
	PacketType      uint8
	Dup             byte
	QoS1            byte
	QoS2            byte
	Retain          byte
	RemainingLength uint
}

var ErrBytesLength = errors.New("fixed header bytes length is between 2 and 5")

// ToFixedHeader converts bytes into a FixedHeader structure.
func ToFixedHeader(bs []byte) (FixedHeader, error) {
	if len(bs) < 2 || 5 < len(bs) {
		return FixedHeader{}, ErrBytesLength
	}
	packetType := bs[0] >> 4
	dup := refbit(bs[0], 3)
	qos1 := refbit(bs[0], 2)
	qos2 := refbit(bs[0], 1)
	retain := refbit(bs[0], 0)
	result := FixedHeader{packetType, dup, qos1, qos2, retain, decodeRemainingLength(bs[1:])}
	return result, nil
}

func refbit(i byte, b uint) byte {
	return (i >> b) & 1
}

func decodeRemainingLength(bs []byte) uint {
	multiplier := uint(1)
	var value uint
	for i := uint(0); i < 8; i++ {
		digit := bs[i]
		value = value + uint(digit&127)*multiplier
		multiplier = multiplier * 128
		if (digit & 128) == 0 {
			break
		}
	}
	return value
}
