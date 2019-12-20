package packet

import (
	"bufio"
	"io"
)

type MQTTReader struct {
	byte1  *byte
	conn   *MQTTConn
	reader *bufio.Reader
}

func NewMQTTReader(conn *MQTTConn) *MQTTReader {
	return &MQTTReader{conn: conn}
}

func (d *MQTTReader) Read(p []byte) (n int, err error) {
	if d.reader == nil {
		r, err := d.conn.nextReader()
		if err != nil {
			return 0, err
		}
		d.reader = r
	}
	return d.reader.Read(p)
}

func (d *MQTTReader) ReadByte() (uint8, error) {
	if d.conn.isWebsocket() {
		if d.reader == nil {
			r, err := d.conn.nextReader()
			if err != nil {
				return 0, err
			}
			d.reader = bufio.NewReader(r)
		}
		for {
			byte1, err := d.reader.ReadByte()
			if err != nil {
				if err == io.EOF {
					r, err := d.conn.nextReader()
					if err != nil {
						return 0, err
					}
					d.reader = bufio.NewReader(r)
					continue
				}
				return 0, err
			}
			return byte1, nil
		}
	}
	if d.reader == nil {
		r, err := d.conn.nextReader()
		if err != nil {
			return 0, err
		}
		d.reader = r
	}
	return d.reader.ReadByte()
}

func (d *MQTTReader) ReadPacketType() (uint8, error) {
	if d.byte1 == nil {
		byte1, err := d.ReadByte()
		if err != nil {
			if err != io.EOF {
				return 0, err
			}
			if byte1 == 0 {
				return 0, io.EOF
			}
		}
		d.byte1 = &byte1
	}
	return *d.byte1 >> 4, nil
}
