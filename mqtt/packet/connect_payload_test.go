package packet_test

import (
	"testing"

	"github.com/bati11/oreno-mqtt/mqtt/packet"
)

func TestToConnectPayload(t *testing.T) {
	in := []byte{0x00, 0x04, 'h', 'o', 'g', 'e'}
	connectPayload := packet.ToConnectPayload(in)
	if connectPayload.ClientID != "hoge" {
		t.Errorf("connectPayload.ClientID: got %v, want %v", connectPayload.ClientID, "hoge")
	}
}
