package packet_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/bati11/oreno-mqtt/mqtt/packet"
)

func sampleVariableHeaderBytes() []byte {
	vb1 := byte(0x00)  // 00000000
	vb2 := byte(0x04)  // 00000100
	vb3 := byte('M')   // 01001101
	vb4 := byte('Q')   // 01010001
	vb5 := byte('T')   // 01010100
	vb6 := byte('T')   // 01010100
	vb7 := byte(0x04)  // 00000100
	vb8 := byte(0x02)  // 00000010
	vb9 := byte(0x00)  // 00000000
	vb10 := byte(0x3C) // 00031100
	return []byte{vb1, vb2, vb3, vb4, vb5, vb6, vb7, vb8, vb9, vb10}
}

func TestToVariableHeader(t *testing.T) {
	variableHeaderBytes := sampleVariableHeaderBytes()
	fixedHeader := packet.FixedHeader{
		PacketType:      packet.CONNECT,
		RemainingLength: uint(len(variableHeaderBytes)),
	}
	variableHeader, _, err := packet.ToConnectVariableHeader(fixedHeader, variableHeaderBytes)
	if err != nil {
		t.Errorf("ToConnectVariableHeader returns err: %v", err)
	}

	if variableHeader.ProtocolName != "MQTT" {
		t.Errorf("ProtocolName: got %v, want %v", variableHeader.ProtocolName, "MQTT")
	}
	if variableHeader.ProtocolLevel != 4 {
		t.Errorf("ProtocolLevel: got %v, want %v", variableHeader.ProtocolLevel, 4)
	}
	if variableHeader.ConnectFlags.CleanSession != true {
		t.Errorf("ConnectFlags.CleanSession: got %v, want %v", variableHeader.ConnectFlags.CleanSession, true)
	}
	if variableHeader.ConnectFlags.WillFlag != false {
		t.Errorf("ConnectFlags.WillFlag: got %v, want %v", variableHeader.ConnectFlags.WillFlag, false)
	}
	if variableHeader.ConnectFlags.WillQoS != 0 {
		t.Errorf("ConnectFlags.WillQoS: got %v, want %v", variableHeader.ConnectFlags.WillQoS, 0)
	}
	if variableHeader.ConnectFlags.WillRetain != false {
		t.Errorf("ConnectFlags.WillRetain: got %v, want %v", variableHeader.ConnectFlags.WillRetain, false)
	}
	if variableHeader.ConnectFlags.PasswordFlag != false {
		t.Errorf("ConnectFlags.PasswordFlag: got %v, want %v", variableHeader.ConnectFlags.PasswordFlag, false)
	}
	if variableHeader.ConnectFlags.UserNameFlag != false {
		t.Errorf("ConnectFlags.UserNameFlag: got %v, want %v", variableHeader.ConnectFlags.UserNameFlag, false)
	}
	if variableHeader.KeepAlive != 60 {
		t.Errorf("KeepAlive: got %v, want %v", variableHeader.KeepAlive, 60)
	}
}

func TestToVariableHeaderInvalidPacketType(t *testing.T) {
	variableHeaderBytes := sampleVariableHeaderBytes()
	fixedHeader := packet.FixedHeader{
		PacketType:      packet.PUBLISH,
		RemainingLength: uint(len(variableHeaderBytes)),
	}
	_, _, err := packet.ToConnectVariableHeader(fixedHeader, variableHeaderBytes)
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
			variableHeaderBytes := tt.in
			vb7 := byte(0x04)  // 00000100
			vb8 := byte(0xCE)  // 11001110
			vb9 := byte(0x00)  // 00000000
			vb10 := byte(0x0A) // 00001010
			variableHeaderBytes = append(variableHeaderBytes, vb7, vb8, vb9, vb10)
			fixedHeader := packet.FixedHeader{
				PacketType:      packet.CONNECT,
				RemainingLength: uint(len(variableHeaderBytes)),
			}
			_, _, err := packet.ToConnectVariableHeader(fixedHeader, variableHeaderBytes)
			if err == nil {
				t.Errorf("ToConnectVariableHeader() returns err: got nil, want error")
			}
		})
	}
}

func TestProtocolLevelInConnect(t *testing.T) {
	var cases = []struct {
		in      byte
		wantErr bool
	}{
		{3, true},
		{4, false},
		{5, true},
	}
	for _, tt := range cases {
		t.Run(string(tt.in), func(t *testing.T) {
			variableHeaderBytes := sampleVariableHeaderBytes()
			variableHeaderBytes[6] = tt.in
			fixedHeader := packet.FixedHeader{
				PacketType:      packet.CONNECT,
				RemainingLength: uint(len(variableHeaderBytes)),
			}
			_, _, err := packet.ToConnectVariableHeader(fixedHeader, variableHeaderBytes)
			if tt.wantErr && (err == nil) {
				t.Errorf("ToConnectVariableHeader() should returns err: but got nil")
			}
			if !tt.wantErr && (err != nil) {
				t.Errorf("ToConnectVariableHeader() should not returns err: but got %v", err)
			}
		})
	}
}

func TestReservedValueInConnectFlagsInConnect(t *testing.T) {
	var cases = []struct {
		in      byte
		wantErr bool
	}{
		{0x02, false},
		{0x03, true},
	}
	for _, tt := range cases {
		t.Run(fmt.Sprintf("%v", tt.in), func(t *testing.T) {
			variableHeaderBytes := sampleVariableHeaderBytes()
			variableHeaderBytes[7] = tt.in
			fixedHeader := packet.FixedHeader{
				PacketType:      packet.CONNECT,
				RemainingLength: uint(len(variableHeaderBytes)),
			}
			_, _, err := packet.ToConnectVariableHeader(fixedHeader, variableHeaderBytes)
			if tt.wantErr && (err == nil) {
				t.Errorf("ToConnectVariableHeader() should returns err: but got nil")
			}
			if !tt.wantErr && (err != nil) {
				t.Errorf("ToConnectVariableHeader() should not returns err: but got %v", err)
			}
		})
	}
}

func TestToStructRevertToByteArray(t *testing.T) {
	variableHeaderBytes := sampleVariableHeaderBytes()

	fixHeader := packet.FixedHeader{
		PacketType:      packet.CONNECT,
		RemainingLength: uint(len(variableHeaderBytes)),
	}

	s, _, err := packet.ToConnectVariableHeader(fixHeader, variableHeaderBytes)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if !reflect.DeepEqual(s.ToBytes(), variableHeaderBytes) {
		t.Fatalf("want %v, got %v", variableHeaderBytes, s.ToBytes())
	}
}
