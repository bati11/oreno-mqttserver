package packet

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"
)

func TestToSubscribePayload(t *testing.T) {
	tests := []struct {
		name    string
		in      []byte
		want    SubscribePayload
		wantErr bool
	}{
		{
			"a/b:QoS1",
			[]byte{0x00, 0x03, 0x61, 0x2F, 0x62, 0x01}, // a/b, QoS1
			SubscribePayload{TopicFilter: "a/b", QoS: QoS1},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToSubscribePayload(bufio.NewReader(bytes.NewBuffer(tt.in)))
			if (err != nil) != tt.wantErr {
				t.Errorf("ToSubscribePayload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToSubscribePayload() = %v, want %v", got, tt.want)
			}
		})
	}
}
