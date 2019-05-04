package packet_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/bati11/oreno-mqtt/mqtt/packet"
)

func TestMQTTReader_ReadPublishVariableHeader(t *testing.T) {
	type args struct {
		r *packet.MQTTReader
	}
	tests := []struct {
		name    string
		args    args
		want    *packet.PublishVariableHeader
		wantErr bool
	}{
		{
			name: "a/b",
			args: args{
				packet.NewMQTTReader(bytes.NewBuffer([]byte{
					0x00,             // Length MSB
					0x03,             // Length LSB
					0x61, 0x2F, 0x62, // a/b
				})),
			},
			want:    &packet.PublishVariableHeader{LengthMSB: 0, LengthLSB: 3, TopicName: "a/b", PacketIdentifier: nil},
			wantErr: false,
		},
		{
			name: "256文字",
			args: args{
				packet.NewMQTTReader(bytes.NewBuffer([]byte{
					0x01,
					0x00,
					'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5',
				})),
			},
			want:    &packet.PublishVariableHeader{LengthMSB: 1, LengthLSB: 0, TopicName: "0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := packet.ExportReadPublishVariableHeader(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExportReadPublishVariableHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExportReadPublishVariableHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPublishVariableHeader_Length(t *testing.T) {
	variableHeaderBytes := []byte{
		0x00,             // Length LSB
		0x03,             // Length MSB
		0x61, 0x2F, 0x62, // a/b
	}
	want := uint(len(variableHeaderBytes))
	r := packet.NewMQTTReader(bytes.NewBuffer(variableHeaderBytes))
	variableHeader, err := packet.ExportReadPublishVariableHeader(r)
	if err != nil {
		t.Errorf("ExportReadPublishVariableHeader() error = %v, wantErr %v", err, false)
		return
	}
	got := variableHeader.Length()
	if got != want {
		t.Errorf("ExportReadPublishVariableHeader.Length() = %v, want %v", got, want)
	}
}
