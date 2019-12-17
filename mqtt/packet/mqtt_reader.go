package packet

import (
	"bufio"
	"io"
)

type MQTTReader struct {
	byte1 *byte
	r     *bufio.Reader
}

func NewMQTTReader(r *bufio.Reader) *MQTTReader {
	return &MQTTReader{r: r}
}

func (d *MQTTReader) ReadPacketType() (uint8, error) {
	if d.byte1 == nil {
		byte1, err := d.r.ReadByte()
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
