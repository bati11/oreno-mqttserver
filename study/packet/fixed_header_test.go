package packet_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/bati11/oreno-mqtt/study/packet"
)

func TestToFixedHeader(t *testing.T) {
	type args struct {
		bs []byte
	}
	tests := []struct {
		args args
		want packet.FixedHeader
	}{
		{
			args{[]byte{
				0x00, // 0000 0 00 0
				0x00, // 0
			}},
			packet.FixedHeader{PacketType: 0, Dup: 0, QoS1: 0, QoS2: 0, Retain: 0, RemainingLength: 0},
		},
		{
			args{[]byte{
				0x1B, // 0001 1 01 1
				0x7F, // 127
			}},
			packet.FixedHeader{PacketType: 1, Dup: 1, QoS1: 0, QoS2: 1, Retain: 1, RemainingLength: 127},
		},
		{
			args{[]byte{
				0x24,       // 0002 0 10 0
				0x80, 0x01, //128
			}},
			packet.FixedHeader{PacketType: 2, Dup: 0, QoS1: 1, QoS2: 0, Retain: 0, RemainingLength: 128},
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%#v", tt.args.bs), func(t *testing.T) {
			if got := packet.ToFixedHeader(tt.args.bs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToFixedHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}
