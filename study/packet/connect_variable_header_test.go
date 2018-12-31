package packet_test

import (
	"reflect"
	"testing"

	"github.com/bati11/oreno-mqtt/study/packet"
)

func TestToConnectVariableHeader(t *testing.T) {
	type args struct {
		fixedHeader packet.FixedHeader
		bs          []byte
	}
	tests := []struct {
		name    string
		args    args
		want    packet.ConnectVariableHeader
		wantErr bool
	}{
		{
			name: "仕様書のexample",
			args: args{
				fixedHeader: packet.FixedHeader{PacketType: 1},
				bs: []byte{
					0x00, 0x04, 'M', 'Q', 'T', 'T', // Protocol Name
					0x04,       // Protocol Level
					0xCE,       // Connect Flags
					0x00, 0x0A, // Keep Alive
				},
			},
			want: packet.ConnectVariableHeader{
				ProtocolName:  "MQTT",
				ProtocolLevel: 4,
				ConnectFlags:  packet.ConnectFlags{UserNameFlag: true, PasswordFlag: true, WillRetain: false, WillQoS: 1, WillFlag: true, CleanSession: true},
				KeepAlive:     10,
			},
			wantErr: false,
		},
		{
			name: "固定ヘッダーのPacketTypeが1ではない",
			args: args{
				fixedHeader: packet.FixedHeader{PacketType: 2},
				bs: []byte{
					0x00, 0x04, 'M', 'Q', 'T', 'T', // Protocol Name
					0x04,       // Protocol Level
					0xCE,       // Connect Flags
					0x00, 0x0A, // Keep Alive
				},
			},
			want:    packet.ConnectVariableHeader{},
			wantErr: true,
		},
		{
			name: "Protocol Nameが不正",
			args: args{
				fixedHeader: packet.FixedHeader{PacketType: 1},
				bs: []byte{
					0x00, 0x04, 'M', 'Q', 'T', 't', // Protocol Name
					0x04,       // Protocol Level
					0xCE,       // Connect Flags
					0x00, 0x0A, // Keep Alive
				},
			},
			want:    packet.ConnectVariableHeader{},
			wantErr: true,
		},
		{
			name: "Protocol Levelが不正",
			args: args{
				fixedHeader: packet.FixedHeader{PacketType: 1},
				bs: []byte{
					0x00, 0x04, 'M', 'Q', 'T', 'T', // Protocol Name
					0x03,       // Protocol Level
					0xCE,       // Connect Flags
					0x00, 0x0A, // Keep Alive
				},
			},
			want:    packet.ConnectVariableHeader{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := packet.ToConnectVariableHeader(tt.args.fixedHeader, tt.args.bs)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToConnectVariableHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToConnectVariableHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}
