package packet_test

import (
	"bytes"
	"net"
	"reflect"
	"testing"

	"github.com/bati11/oreno-mqtt/mqtt/packet"
)

type DummyConn struct {
	net.Conn
	bs *bytes.Buffer
}

func (c *DummyConn) Read(b []byte) (n int, err error) {
	return c.bs.Read(b)
}

func TestMQTTReader_ReadConnectPayload(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *packet.ConnectPayload
		wantErr bool
	}{
		{
			name:    "ClientIDが1文字",
			args:    args{[]byte{0x00, 0x01, 'a'}},
			want:    &packet.ConnectPayload{ClientID: "a"},
			wantErr: false,
		},
		{
			name:    "ペイロードが0byte",
			args:    args{[]byte{}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "ClientIDが23文字を超える",
			args:    args{[]byte{0x00, 0x18, '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', 'a', 'b', 'c', 'd'}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "使えない文字がある",
			args:    args{[]byte{0x00, 0x02, '1', '%'}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "指定された長さよりも実際に取得できたClientIDが短い",
			args:    args{[]byte{0x00, 0x03, '1', '2'}},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := packet.NewMQTTConn(&DummyConn{bs: bytes.NewBuffer(tt.args.b)})
			r := packet.NewMQTTReader(c)
			got, err := packet.ExportReadConnectPayload(r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExportReadConnectPayload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExportReadConnectPayload() = %v, want %v", got, tt.want)
			}
		})
	}
}
