package packet_test

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"

	"github.com/bati11/oreno-mqtt/mqtt/packet"
)

func TestToFixedHeader(t *testing.T) {
	type args struct {
		r *bufio.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    packet.FixedHeader
		wantErr bool
	}{
		{
			name: "[0x00,0x00]",
			args: args{bufio.NewReader(bytes.NewBuffer([]byte{
				0x30, // 0003 0 00 0
				0x00, // 0
			}))},
			want:    packet.PublishFixedHeader{PacketType: 3, Dup: 0, QoS1: 0, QoS2: 0, Retain: 0, RemainingLength: 0},
			wantErr: false,
		},
		{
			name: "[0x1b,0x7F]",
			args: args{bufio.NewReader(bytes.NewBuffer([]byte{
				0x3B, // 0003 1 01 1
				0x7F, // 127
			}))},
			want:    packet.PublishFixedHeader{PacketType: 3, Dup: 1, QoS1: 0, QoS2: 1, Retain: 1, RemainingLength: 127},
			wantErr: false,
		},
		{
			name: "[0x24,0x80,0x01]",
			args: args{bufio.NewReader(bytes.NewBuffer([]byte{
				0x34,       // 0003 0 10 0
				0x80, 0x01, //128
			}))},
			want:    packet.PublishFixedHeader{PacketType: 3, Dup: 0, QoS1: 1, QoS2: 0, Retain: 0, RemainingLength: 128},
			wantErr: false,
		},
		{
			name:    "[]",
			args:    args{bufio.NewReader(bytes.NewBuffer(nil))},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "[0x24]",
			args:    args{bufio.NewReader(bytes.NewBuffer([]byte{0x24}))},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := packet.ToFixedHeader(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToFixedHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToFixedHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFixedHeader_ToBytes(t *testing.T) {
	type fields struct {
		PacketType      byte
		Dup             byte
		QoS1            byte
		QoS2            byte
		Retain          byte
		RemainingLength uint
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			"Remining Length = 0",
			fields{PacketType: 1, Dup: 0, QoS1: 0, QoS2: 0, Retain: 0, RemainingLength: 0},
			[]byte{
				0x10,
				0x00,
			},
		},
		{
			"Remining Length = 127",
			fields{PacketType: 1, Dup: 0, QoS1: 0, QoS2: 0, Retain: 0, RemainingLength: 127},
			[]byte{
				0x10,
				0x7F,
			},
		},
		{
			"Remining Length = 128",
			fields{PacketType: 1, Dup: 0, QoS1: 0, QoS2: 0, Retain: 0, RemainingLength: 128},
			[]byte{
				0x10,
				0x80, 0x01,
			},
		},
		{
			"Remining Length = 16383",
			fields{PacketType: 1, Dup: 0, QoS1: 0, QoS2: 0, Retain: 0, RemainingLength: 16383},
			[]byte{
				0x10,
				0xFF, 0x7F,
			},
		},
		{
			"Remining Length = 16384",
			fields{PacketType: 1, Dup: 0, QoS1: 0, QoS2: 0, Retain: 0, RemainingLength: 16384},
			[]byte{
				0x10,
				0x80, 0x80, 0x01,
			},
		},
		{
			"Remining Length = 2097151",
			fields{PacketType: 1, Dup: 0, QoS1: 0, QoS2: 0, Retain: 0, RemainingLength: 2097151},
			[]byte{
				0x10,
				0xFF, 0xFF, 0x7F,
			},
		},
		{
			"Remining Length = 2097152",
			fields{PacketType: 1, Dup: 0, QoS1: 0, QoS2: 0, Retain: 0, RemainingLength: 2097152},
			[]byte{
				0x10,
				0x80, 0x80, 0x80, 0x01,
			},
		},
		{
			"Remining Length = 268435455",
			fields{PacketType: 1, Dup: 0, QoS1: 0, QoS2: 0, Retain: 0, RemainingLength: 268435455},
			[]byte{
				0x10,
				0xFF, 0xFF, 0xFF, 0x7F,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := packet.PublishFixedHeader{
				PacketType:      tt.fields.PacketType,
				Dup:             tt.fields.Dup,
				QoS1:            tt.fields.QoS1,
				QoS2:            tt.fields.QoS2,
				Retain:          tt.fields.Retain,
				RemainingLength: tt.fields.RemainingLength,
			}
			if got := h.ToBytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultFixedHeader.ToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
