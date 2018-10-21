package packet

import (
	"reflect"
	"testing"
)

func createUint16(x uint16) *uint16 {
	return &x
}

func TestToPublishVariableHeader(t *testing.T) {

	type args struct {
		fixedHeader FixedHeader
		bs          []byte
	}
	tests := []struct {
		name    string
		args    args
		want    PublishVariableHeader
		want1   []byte
		wantErr bool
	}{
		{
			"QoS0",
			args{
				FixedHeader{PacketType: PUBLISH, QoS: QoS0},
				[]byte{0x00, 0x03, 0x61, 0x2F, 0x62, 0x00, 0x0A}, // a/b
			},
			PublishVariableHeader{TopicName: "a/b"},
			[]byte{0x00, 0x0A},
			false,
		},
		{
			"QoS1",
			args{
				FixedHeader{PacketType: PUBLISH, QoS: QoS1},
				[]byte{0x00, 0x03, 0x61, 0x2F, 0x62, 0x00, 0x0A}, // a/b, 10
			},
			PublishVariableHeader{TopicName: "a/b", PacketIdentifier: createUint16(10)},
			[]byte{},
			false,
		},
		{
			"QoS2",
			args{
				FixedHeader{PacketType: PUBLISH, QoS: QoS2},
				[]byte{0x00, 0x03, 0x61, 0x2F, 0x62, 0x00, 0x0A}, // a/b, 10
			},
			PublishVariableHeader{TopicName: "a/b", PacketIdentifier: createUint16(10)},
			[]byte{},
			false,
		},
		{
			"len = 0",
			args{
				FixedHeader{PacketType: PUBLISH, QoS: QoS0},
				[]byte{},
			},
			PublishVariableHeader{}, nil,
			true,
		},
		{
			"len = 2",
			args{
				FixedHeader{PacketType: PUBLISH, QoS: QoS0},
				[]byte{0x00, 0x03},
			},
			PublishVariableHeader{}, nil,
			true,
		},
		{
			"Length LSB = 0",
			args{
				FixedHeader{PacketType: PUBLISH, QoS: QoS0},
				[]byte{0x00, 0x00}, // MSB=0, LSB=0
			},
			PublishVariableHeader{}, nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ToPublishVariableHeader(tt.args.fixedHeader, tt.args.bs)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToPublishVariableHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToPublishVariableHeader() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ToPublishVariableHeader() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
