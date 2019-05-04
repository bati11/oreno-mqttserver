package packet

import (
	"bufio"
	"io"
)

type MQTTReader struct {
	byte1 *byte
	r     *bufio.Reader
}

func NewMQTTReader(r io.Reader) *MQTTReader {
	bufr := bufio.NewReader(r)
	return &MQTTReader{r: bufr}
}

func (d *MQTTReader) ReadPacketType() (uint8, error) {
	if d.byte1 == nil {
		byte1, err := d.r.ReadByte()
		if err != nil {
			return 0, err
		}
		d.byte1 = &byte1
	}
	return *d.byte1 >> 4, nil
}
