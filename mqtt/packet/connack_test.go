package packet_test

import (
	"bytes"
	"testing"

	"github.com/bati11/oreno-mqtt/mqtt/packet"
)

func TestConnackToBytes(t *testing.T) {
	cases := []struct {
		name string
		in   packet.Connack
		want []byte
	}{
		{
			"Connection Accepted",
			packet.NewConnackForAccepted(),
			[]byte{0x20, 0x02, 0x00, 0x00},
		},
		{
			"Connection Refused, unacceptable protocol version",
			packet.NewConnackForRefusedByUnacceptableProtocolVersion(),
			[]byte{0x20, 0x02, 0x00, 0x01},
		},
		{
			"Connection Refused, identifier rejected",
			packet.NewConnackForRefusedByIdentifierRejected(),
			[]byte{0x20, 0x02, 0x00, 0x02},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.in.ToBytes()
			if !bytes.Equal(result, tt.want) {
				t.Fatalf("want %v, got %v", tt.want, result)
			}
		})
	}
}
