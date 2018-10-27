package handler

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/bati11/oreno-mqtt/mqtt/packet"
)

func HandlePublish(fixedHeader packet.FixedHeader, r *bufio.Reader) (*packet.MQTTPacket, error) {
	bs := make([]byte, 6)
	_, err := io.ReadFull(r, bs)
	if err != nil {
		return nil, err
	}
	variableHeader, _, err := packet.ToPublishVariableHeader(fixedHeader, bs)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%v\n", variableHeader)

	lReader := io.LimitReader(r, int64(fixedHeader.RemainingLength-6))
	io.Copy(os.Stdout, lReader)
	fmt.Println()
	//payload := packet.ToPublishPayload(remains)
	//fmt.Printf("%v\n", payload)
	return nil, nil
}
