package packet_test

import (
	"testing"

	"github.com/bati11/oreno-mqtt/mqtt/packet"
)

func TestToFixedHeader(t *testing.T) {
	b1 := byte(0x12) // 0001 0 01 0
	b2 := byte(0x00) // 00000000
	in := [2]byte{b1, b2}

	result := packet.ToFixedHeader(in)
	if result.PacketType != 1 {
		t.Errorf("PacketType: got %q, want %q", result.PacketType, 1)
	}
	if result.Dup != 0 {
		t.Errorf("Dup: got %q, want %q", result.Dup, 0)
	}
	if result.QoS1 != 0 {
		t.Errorf("QoS1: got %q, want %q", result.QoS1, 0)
	}
	if result.QoS2 != 1 {
		t.Errorf("QoS2: got %q, want %q", result.QoS2, 1)
	}
	if result.Retain != 0 {
		t.Errorf("Retain: got %q, want %q", result.Retain, 0)
	}
	if result.RemainingLength != 0 {
		t.Errorf("RemainingLength: got %q, want %q", result.RemainingLength, 0)
	}
}
