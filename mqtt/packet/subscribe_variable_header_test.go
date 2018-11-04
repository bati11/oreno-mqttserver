package packet

import (
	"fmt"
	"reflect"
	"testing"
)

func TestToSubscribeVariableHeader(t *testing.T) {
	tests := []struct {
		in      []byte
		want    SubscribeVariableHeader
		wantErr bool
	}{
		{
			[]byte{0x00, 0x0A},
			SubscribeVariableHeader{PacketIdentifier: 10},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%X", tt.in), func(t *testing.T) {
			got, err := ToSubscribeVariableHeader(tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToSubscribeVariableHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToSubscribeVariableHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}
