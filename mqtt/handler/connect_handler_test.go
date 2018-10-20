package handler_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/bati11/oreno-mqtt/mqtt/handler"
	"github.com/bati11/oreno-mqtt/mqtt/packet"
)

func TestHandleConnect(t *testing.T) {
	type In struct {
		fixedHeader    packet.FixedHeader
		variableHeader packet.ConnectVariableHeader
		connectPayload packet.ConnectPayload
	}
	connectFixedHeader := packet.FixedHeader{PacketType: packet.CONNECT}
	cases := []struct {
		name      string
		in        In
		want      packet.Connack
		wantError bool
		err       error
	}{
		{
			name: "Accepted",
			in: In{
				connectFixedHeader,
				packet.ConnectVariableHeader{ProtocolName: "MQTT", ProtocolLevel: 4, ConnectFlags: packet.ConnectFlags{CleanSession: true}},
				packet.ConnectPayload{"hogehoge"},
			},
			want:      packet.NewConnackForAccepted(),
			wantError: false,
		},
		{
			name: "UnacceptableProtocolVersion",
			in: In{
				connectFixedHeader,
				packet.ConnectVariableHeader{ProtocolName: "MQTT", ProtocolLevel: 3, ConnectFlags: packet.ConnectFlags{CleanSession: true}},
				packet.ConnectPayload{"hogehoge"},
			},
			want:      packet.NewConnackForRefusedByUnacceptableProtocolVersion(),
			wantError: false,
		},
		{
			name: "IdentifierRejected",
			in: In{
				connectFixedHeader,
				packet.ConnectVariableHeader{ProtocolName: "MQTT", ProtocolLevel: 4, ConnectFlags: packet.ConnectFlags{CleanSession: true}},
				packet.ConnectPayload{"あああ"},
			},
			want:      packet.NewConnackForRefusedByIdentifierRejected(),
			wantError: false,
		},
		{
			name: "not MQTT",
			in: In{
				connectFixedHeader,
				packet.ConnectVariableHeader{ProtocolName: "aaa", ProtocolLevel: 4, ConnectFlags: packet.ConnectFlags{CleanSession: true}},
				packet.ConnectPayload{"hogehoge"},
			},
			want:      packet.NewConnackForRefusedByIdentifierRejected(),
			wantError: true,
			err:       errors.New("aaa"),
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			var bs []byte
			bs = append(bs, tt.in.variableHeader.ToBytes()...)
			bs = append(bs, tt.in.connectPayload.ToBytes()...)
			tt.in.fixedHeader.RemainingLength = uint(len(bs))

			result, err := handler.HandleConnect(tt.in.fixedHeader, bs)
			if !tt.wantError && (err != nil) {
				t.Fatalf("want no err, but %#v", err)
			}
			if tt.wantError && !(reflect.DeepEqual(err, tt.err)) {
				t.Fatalf("want %v, but %#v", tt.err, err)
			}
			if !reflect.DeepEqual(result, tt.want) {
				t.Fatalf("want %+v, got %+v", tt.want, result)
			}
		})
	}
}
