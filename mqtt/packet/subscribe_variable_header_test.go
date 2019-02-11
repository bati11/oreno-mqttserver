package packet

import (
	"bufio"
	"reflect"
	"testing"
)

func TestToSubscribeVariableHeader(t *testing.T) {
	type args struct {
		fixedHeader FixedHeader
		r           *bufio.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    SubscribeVariableHeader
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToSubscribeVariableHeader(tt.args.fixedHeader, tt.args.r)
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
