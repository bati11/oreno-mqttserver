package packet_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/bati11/oreno-mqtt/mqtt/packet"
)

func TestMQTTReader_ReadConnectVariableHeader(t *testing.T) {
	type args struct {
		r *packet.MQTTReader
	}
	tests := []struct {
		name    string
		args    args
		want    *packet.ConnectVariableHeader
		wantErr bool
	}{
		{
			name: "仕様書のexample",
			args: args{
				r: packet.NewMQTTReader(bytes.NewBuffer([]byte{
					0x00, 0x04, 'M', 'Q', 'T', 'T', // Protocol Name
					0x04,       // Protocol Level
					0xCE,       // Connect Flags
					0x00, 0x0A, // Keep Alive
				})),
			},
			want: &packet.ConnectVariableHeader{
				ProtocolName:  "MQTT",
				ProtocolLevel: 4,
				ConnectFlags:  packet.ConnectFlags{UserNameFlag: true, PasswordFlag: true, WillRetain: false, WillQoS: 1, WillFlag: true, CleanSession: true},
				KeepAlive:     10,
			},
			wantErr: false,
		},
		{
			name: "Protocol Nameが不正",
			args: args{
				r: packet.NewMQTTReader(bytes.NewReader([]byte{
					0x00, 0x04, 'M', 'Q', 'T', 't', // Protocol Name
					0x04,       // Protocol Level
					0xCE,       // Connect Flags
					0x00, 0x0A, // Keep Alive
				})),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Protocol Levelが不正",
			args: args{
				r: packet.NewMQTTReader(bytes.NewReader([]byte{
					0x00, 0x04, 'M', 'Q', 'T', 'T', // Protocol Name
					0x03,       // Protocol Level
					0xCE,       // Connect Flags
					0x00, 0x0A, // Keep Alive
				})),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := packet.ExportReadVariableConnectHeader(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExportReadVariableConnectHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExportReadVariableConnectHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}
