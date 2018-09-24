package packet_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/bati11/oreno-mqtt/mqtt/packet"
)

func TestToVariableHeader(t *testing.T) {
	fb1 := byte(0x10) // 0001 0 00 0
	fb2 := byte(0x01) // 00000000
	fixedHeaderBytes := []byte{fb1, fb2}

	vb1 := byte(0x00)  // 00000000
	vb2 := byte(0x04)  // 00000100
	vb3 := byte('M')   // 01001101
	vb4 := byte('Q')   // 01010001
	vb5 := byte('T')   // 01010100
	vb6 := byte('T')   // 01010100
	vb7 := byte(0x04)  // 00000100
	vb8 := byte(0xCE)  // 11001110
	vb9 := byte(0x00)  // 00000000
	vb10 := byte(0x0A) // 00001010
	variableHeaderBytes := []byte{vb1, vb2, vb3, vb4, vb5, vb6, vb7, vb8, vb9, vb10}

	in := append(fixedHeaderBytes, variableHeaderBytes...)

	fixedHeader, remains, err := packet.ToFixedHeader(in)
	if err != nil {
		t.Errorf("ToFixedHeader() returns err: %v", err)
	}
	if !bytes.Equal(remains, variableHeaderBytes) {
		t.Errorf("remains: got %v, want %v", remains, variableHeaderBytes)
	}

	variableHeader, err := packet.ToConnectVariableHeader(fixedHeader, remains)
	if err != nil {
		t.Errorf("ToConnectVariableHeader() returns err: %v", err)
	}

	if variableHeader.ProtocolName() != "MQTT" {
		t.Errorf("ProtocolName(): got %v, want %v", variableHeader.ProtocolName(), "MQTT")
	}
}

func TestToVariableHeaderInvalidPacketType(t *testing.T) {
	fixedHeaderBytes := []byte{byte(packet.PUBLISH), 0x01}

	vb1 := byte(0x00)  // 00000000
	vb2 := byte(0x04)  // 00000100
	vb3 := byte('M')   // 01001101
	vb4 := byte('Q')   // 01010001
	vb5 := byte('T')   // 01010100
	vb6 := byte('T')   // 01010100
	vb7 := byte(0x04)  // 00000100
	vb8 := byte(0xCE)  // 11001110
	vb9 := byte(0x00)  // 00000000
	vb10 := byte(0x0A) // 00001010
	variableHeaderBytes := []byte{vb1, vb2, vb3, vb4, vb5, vb6, vb7, vb8, vb9, vb10}

	in := append(fixedHeaderBytes, variableHeaderBytes...)

	fixedHeader, remains, err := packet.ToFixedHeader(in)
	if err != nil {
		t.Errorf("ToFixedHeader() returns err: %v", err)
	}
	if !bytes.Equal(remains, variableHeaderBytes) {
		t.Errorf("remains: got %v, want %v", remains, variableHeaderBytes)
	}

	_, err = packet.ToConnectVariableHeader(fixedHeader, remains)
	if err == nil {
		t.Errorf("ToConnectVariableHeader() returns err: got nil, want err")
	}
}

func TestProtocolNameInConnect(t *testing.T) {
	var cases = []struct {
		in []byte
	}{
		{[]byte{0, 4, 'm', 'Q', 'T', 'T'}},
		{[]byte{1, 4, 'M', 'Q', 'T', 'T'}},
		{[]byte{0, 5, 'M', 'Q', 'T', 'T'}},
	}
	for _, tt := range cases {
		t.Run(fmt.Sprintf("%q", tt.in), func(t *testing.T) {
			fixedHeaderBytes := []byte{0x10, 0x01}
			variableHeaderBytes := tt.in
			vb7 := byte(0x04)  // 00000100
			vb8 := byte(0xCE)  // 11001110
			vb9 := byte(0x00)  // 00000000
			vb10 := byte(0x0A) // 00001010
			variableHeaderBytes = append(variableHeaderBytes, vb7, vb8, vb9, vb10)

			in := append(fixedHeaderBytes, variableHeaderBytes...)
			fixedHeader, remains, err := packet.ToFixedHeader(in)
			if err != nil {
				t.Errorf("ToFixedHeader() returns err: %v", err)
			}
			if !bytes.Equal(remains, variableHeaderBytes) {
				t.Errorf("remains: got %v, want %v", remains, variableHeaderBytes)
			}

			_, err = packet.ToConnectVariableHeader(fixedHeader, remains)
			if err == nil {
				t.Errorf("ToConnectVariableHeader() returns err: got nil, want error")
			}
		})
	}
}
