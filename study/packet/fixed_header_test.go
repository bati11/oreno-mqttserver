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
			"Reserved",
			args{[]byte{0x00, 0x00}},
			packet.FixedHeader{0},
		},
		{
			"CONNECT",
			args{[]byte{0x10, 0x00}},
			packet.FixedHeader{1},
		},
		{
			"CONNACK",
			args{[]byte{0x20, 0x00}},
			packet.FixedHeader{2},
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
