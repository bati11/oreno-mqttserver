package packet_test

import (
	"github.com/bati11/oreno-mqtt/study/packet"
	"reflect"
	"testing"
)

func TestToFixedHeader(t *testing.T) {
	type args struct {
		bs []byte
	}
	tests := []struct {
		name string
		args args
		want packet.FixedHeader
	}{
		{
			"Reserved Dup:0 QoS:00 Retain:0",
			args{[]byte{0x00, 0x00}}, // 0000 0 00 0
			packet.FixedHeader{PacketType: 0, Dup: 0, QoS1: 0, QoS2: 0, Retain: 0},
		},
		{
			"CONNECT Dup:1 QoS:01 Retain:1",
			args{[]byte{0x1B, 0x00}}, // 0001 1 01 1
			packet.FixedHeader{PacketType: 1, Dup: 1, QoS1: 0, QoS2: 1, Retain: 1},
		},
		{
			"CONNACK Dup:0 QoS:10 Retain:1",
			args{[]byte{0x24, 0x00}}, // 0002 0 10 0
			packet.FixedHeader{PacketType: 2, Dup: 0, QoS1: 1, QoS2: 0, Retain: 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := packet.ToFixedHeader(tt.args.bs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToFixedHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}
