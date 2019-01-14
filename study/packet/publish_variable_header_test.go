package packet

import (
	"bufio"
	"reflect"
	"testing"
)

func TestToPublishVariableHeader(t *testing.T) {
	type args struct {
		fixedHeader FixedHeader
		r           *bufio.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    PublishVariableHeader
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToPublishVariableHeader(tt.args.fixedHeader, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToPublishVariableHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToPublishVariableHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}
