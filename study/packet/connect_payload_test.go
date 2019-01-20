package packet

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"
)

func TestToConnectPayload(t *testing.T) {
	type args struct {
		r *bufio.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    ConnectPayload
		wantErr bool
	}{
		{
			name:    "ClientIDが1文字",
			args:    args{bufio.NewReader(bytes.NewBuffer([]byte{0x00, 0x01, 'a'}))},
			want:    ConnectPayload{ClientID: "a"},
			wantErr: false,
		},
		{
			name:    "ペイロードが0byte",
			args:    args{bufio.NewReader(bytes.NewBuffer([]byte{}))},
			want:    ConnectPayload{},
			wantErr: true,
		},
		{
			name:    "ClientIDが23文字を超える",
			args:    args{bufio.NewReader(bytes.NewBuffer([]byte{0x00, 0x18, '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', 'a', 'b', 'c', 'd'}))},
			want:    ConnectPayload{},
			wantErr: true,
		},
		{
			name:    "使えない文字がある",
			args:    args{bufio.NewReader(bytes.NewBuffer([]byte{0x00, 0x02, '1', '%'}))},
			want:    ConnectPayload{},
			wantErr: true,
		},
		{
			name:    "指定された長さよりも実際に取得できたClientIDが短い",
			args:    args{bufio.NewReader(bytes.NewBuffer([]byte{0x00, 0x03, '1', '2'}))},
			want:    ConnectPayload{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToConnectPayload(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToConnectPayload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToConnectPayload() = %v, want %v", got, tt.want)
			}
		})
	}
}
