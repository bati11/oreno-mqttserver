package packet_test

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"

	"github.com/bati11/oreno-mqtt/mqtt/packet"
)

func TestToPublishVariableHeader(t *testing.T) {
	type args struct {
		fixedHeader packet.PublishFixedHeader
		r           *bufio.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    packet.PublishVariableHeader
		wantErr bool
	}{
		{
			name: "a/b",
			args: args{
				packet.PublishFixedHeader{PacketType: packet.PUBLISH, RemainingLength: 10},
				bufio.NewReader(bytes.NewBuffer([]byte{
					0x00,             // Length LSB
					0x03,             // Length MSB
					0x61, 0x2F, 0x62, // a/b
				})),
			},
			want:    packet.PublishVariableHeader{TopicName: "a/b", PacketIdentifier: nil},
			wantErr: false,
		},
		{
			name: "256文字",
			args: args{
				packet.PublishFixedHeader{PacketType: packet.PUBLISH, RemainingLength: 10},
				bufio.NewReader(bytes.NewBuffer([]byte{
					0x01,
					0x00,
					'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5',
				})),
			},
			want:    packet.PublishVariableHeader{TopicName: "0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := packet.ToPublishVariableHeader(tt.args.fixedHeader, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToPublishVariableHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToPublishVariableHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPublishVariableHeader_Length(t *testing.T) {
	fixedHeader := packet.PublishFixedHeader{PacketType: packet.PUBLISH, RemainingLength: 10}
	variableHeaderBytes := []byte{
		0x00,             // Length LSB
		0x03,             // Length MSB
		0x61, 0x2F, 0x62, // a/b
	}
	want := uint(len(variableHeaderBytes))
	variableHeader, err := packet.ToPublishVariableHeader(fixedHeader, bufio.NewReader(bytes.NewBuffer(variableHeaderBytes)))
	if err != nil {
		t.Errorf("ToPublishVariableHeader() error = %v, wantErr %v", err, false)
		return
	}
	got := variableHeader.Length()
	if got != want {
		t.Errorf("PublishVariableHeader.Length() = %v, want %v", got, want)
	}
}
