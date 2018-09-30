package packet_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/bati11/oreno-mqtt/mqtt/packet"
)

func TestToConnectPayload(t *testing.T) {
	cases := []struct {
		in        []byte
		want      string
		wantError bool
		err       error
	}{
		{
			[]byte{0x00, 0x04, 'h', 'o', 'g', 'e'}, "hoge",
			false, nil,
		},
		{
			[]byte{0x00, 0x04, 'h', 'o', 'g'}, "hog",
			false, nil,
		},
		{
			[]byte{}, "",
			true, packet.ErrConnectPayloadLength,
		},
		{
			[]byte{0x00, 0x01, '*'}, "",
			true, packet.ErrClientIDFormat,
		},
	}
	for _, tt := range cases {
		t.Run(fmt.Sprintf("%q", tt.in), func(t *testing.T) {
			connectPayload, err := packet.ToConnectPayload(tt.in)
			if !tt.wantError && (err != nil) {
				t.Fatalf("want no err, but %#v", err)
			}
			if tt.wantError && !(reflect.DeepEqual(err, tt.err)) {
				t.Fatalf("want %v, but %#v", tt.err, err)
			}

			if !tt.wantError && connectPayload.ClientID != tt.want {
				t.Fatalf("connectPayload.ClientID: got %v, want %v", connectPayload.ClientID, tt.want)
			}
		})
	}
}
