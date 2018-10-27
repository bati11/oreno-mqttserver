package packet_test

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"

	"github.com/bati11/oreno-mqtt/mqtt/packet"
)

func TestToFixedHeader(t *testing.T) {
	b1 := byte(0x12) // 0001 0 01 0
	b2 := byte(0x01) // 00000000
	b3 := byte(0x00)
	in := []byte{b1, b2, b3}

	result, err := packet.ToFixedHeader(bufio.NewReader(bytes.NewBuffer(in)))
	if err != nil {
		t.Errorf("ToFixedHeader() returns err: %v", err)
	}
	if result.PacketType != 1 {
		t.Errorf("PacketType: got %v, want %v", result.PacketType, 1)
	}
	if result.Dup != false {
		t.Errorf("Dup: got %v, want %v", result.Dup, false)
	}
	if result.QoS != packet.QoS1 {
		t.Errorf("QoS: got %v, want %v", result.QoS, packet.QoS1)
	}
	if result.Retain != false {
		t.Errorf("Retain: got %v, want %v", result.Retain, false)
	}
	if result.RemainingLength != 1 {
		t.Errorf("RemainingLength: got %v, want %v", result.RemainingLength, 1)
	}
}

func TestPacketType(t *testing.T) {
	var cases = []struct {
		in   byte
		want packet.PacketType
	}{
		{0x00, packet.Reserved},
		{0x10, packet.CONNECT},
		{0x20, packet.CONNACK},
		{0x30, packet.PUBLISH},
		{0x40, packet.PUBACK},
		{0x50, packet.PUBREC},
		{0x60, packet.PUBREL},
		{0x70, packet.PUBCOMP},
		{0x80, packet.SUBSCRIBE},
		{0x90, packet.SUBACK},
		{0xA0, packet.UNSUBSCRIBE},
		{0xB0, packet.UNSUBACK},
		{0xC0, packet.PINGREQ},
		{0xD0, packet.PINGRESP},
		{0xE0, packet.DISCONNECT},
		{0xF0, packet.Reserved2},
	}
	for _, tt := range cases {
		t.Run(fmt.Sprintf("%X", tt.in), func(t *testing.T) {
			bs := []byte{tt.in, 0x00}
			result, err := packet.ToFixedHeader(bufio.NewReader(bytes.NewBuffer(bs)))
			if err != nil {
				t.Errorf("ToFixedHeader() returns err: %v", err)
			}
			if result.PacketType != tt.want {
				t.Errorf("PacketType: got %q, want %q", result.PacketType, tt.want)
			}
		})
	}
}

func TestDup(t *testing.T) {
	var cases = []struct {
		in   byte
		want bool
	}{
		{0x10, false},
		{0x18, true},
		{0x28, true},
	}
	for _, tt := range cases {
		t.Run(fmt.Sprintf("%X", tt.in), func(t *testing.T) {
			bs := []byte{tt.in, 0x00}
			result, err := packet.ToFixedHeader(bufio.NewReader(bytes.NewBuffer(bs)))
			if err != nil {
				t.Errorf("ToFixedHeader() returns err: %v", err)
			}
			if result.Dup != tt.want {
				t.Errorf("Dup: got %v, want %v", result.Dup, tt.want)
			}
		})
	}
}

func TestQoS(t *testing.T) {
	var cases = []struct {
		in   byte
		want packet.QoS
	}{
		{0x10, packet.QoS0},
		{0x12, packet.QoS1},
		{0x14, packet.QoS2},
	}
	for _, tt := range cases {
		t.Run(fmt.Sprintf("%X", tt.in), func(t *testing.T) {
			bs := []byte{tt.in, 0x00}
			result, err := packet.ToFixedHeader(bufio.NewReader(bytes.NewBuffer(bs)))
			if err != nil {
				t.Errorf("ToFixedHeader() returns err: %v", err)
			}
			if result.QoS != tt.want {
				t.Errorf("QoS: got %v, want %v", result.QoS, tt.want)
			}
		})
	}
}

func TestRetain(t *testing.T) {
	var cases = []struct {
		in   byte
		want bool
	}{
		{0x10, false},
		{0x11, true},
	}
	for _, tt := range cases {
		t.Run(fmt.Sprintf("%X", tt.in), func(t *testing.T) {
			bs := []byte{tt.in, 0x00}
			result, err := packet.ToFixedHeader(bufio.NewReader(bytes.NewBuffer(bs)))
			if err != nil {
				t.Errorf("ToFixedHeader() returns err: %v", err)
			}
			if result.Retain != tt.want {
				t.Errorf("Retain: got %v, want %v", result.Retain, tt.want)
			}
		})
	}
}

func TestRemainingLength(t *testing.T) {
	var cases = []struct {
		in   []byte
		want uint
	}{
		{[]byte{0x00}, 0},
		{[]byte{0x01}, 1},
		{[]byte{0x7F}, 127},
		{[]byte{0x80, 0x01}, 128},
		{[]byte{0xFF, 0x7F}, 16383},
		{[]byte{0x80, 0x80, 0x01}, 16384},
		{[]byte{0xFF, 0xFF, 0x7F}, 2097151},
		{[]byte{0x80, 0x80, 0x80, 0x01}, 2097152},
		{[]byte{0xFF, 0xFF, 0xFF, 0x7F}, 268435455},
	}
	for _, tt := range cases {
		t.Run(fmt.Sprintf("%X", tt.in), func(t *testing.T) {
			bs := append([]byte{0x10}, tt.in...)
			wantRemains := []byte{0x20, 0x21, 0x22}
			bs = append(bs, wantRemains...)
			result, err := packet.ToFixedHeader(bufio.NewReader(bytes.NewBuffer(bs)))
			if err != nil {
				t.Errorf("ToFixedHeader() returns err: %v", err)
			}
			if result.RemainingLength != tt.want {
				t.Errorf("RemainingLength: got %v, want %v", result.RemainingLength, tt.want)
			}
			//if !bytes.Equal(resultRemains, wantRemains) {
			//	t.Errorf("wantRemains: got %v, want %X", resultRemains, wantRemains)
			//}
		})
	}
}

func TestErrorCase(t *testing.T) {
	var cases = []struct {
		in   []byte
		want error
	}{
		{[]byte{}, packet.ErrBytesLength},
		{[]byte{0x10}, packet.ErrBytesLength},
	}
	for _, tt := range cases {
		t.Run(fmt.Sprintf("%X", tt.in), func(t *testing.T) {
			_, err := packet.ToFixedHeader(bufio.NewReader(bytes.NewBuffer(tt.in)))
			if err == nil {
				t.Errorf("ToFixedHeader: got %v, want %v", err, tt.want)
			}
		})
	}
}

func TestFixedHeaderToBytes(t *testing.T) {
	cases := []struct {
		in        packet.FixedHeader
		want      []byte
		wantError bool
	}{
		{
			packet.FixedHeader{packet.CONNACK, false, packet.QoS0, false, 2},
			[]byte{0x20, 0x02},
			false,
		},
	}
	for _, tt := range cases {
		result := tt.in.ToBytes()
		if !bytes.Equal(result, tt.want) {
			t.Fatalf("got %v, want %v", result, tt.want)
		}
	}
}
