package packet_test

import (
	"fmt"
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

func TestPacketType(t *testing.T) {
	var cases = []struct {
		in   byte
		want uint8
	}{
		{0x10, 1},
		{0x20, 2},
		{0x30, 3},
		{0x40, 4},
		{0x50, 5},
		{0x60, 6},
		{0x70, 7},
		{0x80, 8},
		{0x90, 9},
		{0xA0, 10},
		{0xB0, 11},
		{0xC0, 12},
		{0xD0, 13},
		{0xE0, 14},
		{0xF0, 15},
	}
	for _, tt := range cases {
		t.Run(fmt.Sprintf("%X", tt.in), func(t *testing.T) {
			bs := [2]byte{tt.in, 0x00}
			result := packet.ToFixedHeader(bs)
			if result.PacketType != tt.want {
				t.Errorf("PacketType: got %q, want %q", result.PacketType, tt.want)
			}
		})
	}
}

func TestDup(t *testing.T) {
	var cases = []struct {
		in   byte
		want uint8
	}{
		{0x10, 0},
		{0x18, 1},
		{0x28, 1},
	}
	for _, tt := range cases {
		t.Run(fmt.Sprintf("%X", tt.in), func(t *testing.T) {
			bs := [2]byte{tt.in, 0x00}
			result := packet.ToFixedHeader(bs)
			if result.Dup != tt.want {
				t.Errorf("Dup: got %q, want %q", result.Dup, tt.want)
			}
		})
	}
}

func TestQoS1(t *testing.T) {
	var cases = []struct {
		in   byte
		want uint8
	}{
		{0x10, 0},
		{0x14, 1},
	}
	for _, tt := range cases {
		t.Run(fmt.Sprintf("%X", tt.in), func(t *testing.T) {
			bs := [2]byte{tt.in, 0x00}
			result := packet.ToFixedHeader(bs)
			if result.QoS1 != tt.want {
				t.Errorf("QoS1: got %q, want %q", result.QoS1, tt.want)
			}
		})
	}
}

func TestQoS2(t *testing.T) {
	var cases = []struct {
		in   byte
		want uint8
	}{
		{0x10, 0},
		{0x12, 1},
	}
	for _, tt := range cases {
		t.Run(fmt.Sprintf("%X", tt.in), func(t *testing.T) {
			bs := [2]byte{tt.in, 0x00}
			result := packet.ToFixedHeader(bs)
			if result.QoS2 != tt.want {
				t.Errorf("QoS2: got %q, want %q", result.QoS2, tt.want)
			}
		})
	}
}

func TestRetain(t *testing.T) {
	var cases = []struct {
		in   byte
		want uint8
	}{
		{0x10, 0},
		{0x11, 1},
	}
	for _, tt := range cases {
		t.Run(fmt.Sprintf("%X", tt.in), func(t *testing.T) {
			bs := [2]byte{tt.in, 0x00}
			result := packet.ToFixedHeader(bs)
			if result.Retain != tt.want {
				t.Errorf("Retain: got %q, want %q", result.Retain, tt.want)
			}
		})
	}
}

func TestRemainingLength(t *testing.T) {
	var cases = []struct {
		in   byte
		want uint8
	}{
		{0x00, 0},
		{0x01, 1},
		{0x7F, 127},
	}
	for _, tt := range cases {
		t.Run(fmt.Sprintf("%X", tt.in), func(t *testing.T) {
			bs := [2]byte{0x10, tt.in}
			result := packet.ToFixedHeader(bs)
			if result.RemainingLength != tt.want {
				t.Errorf("RemainingLength: got %v, want %v", result.RemainingLength, tt.want)
			}
		})
	}
}
