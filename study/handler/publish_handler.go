package handler

import (
	"bufio"
	"fmt"
	"io"

	"github.com/bati11/oreno-mqtt/study/packet"
)

func HandlePublish(fixedHeader packet.FixedHeader, r *bufio.Reader) error {
	fmt.Printf("  HandlePublish\n")
	variableHeader, err := packet.ToPublishVariableHeader(fixedHeader, r)
	if err != nil {
		return err
	}
	fmt.Printf("  %#v\n", variableHeader)

	payloadLength := fixedHeader.RemainingLength - variableHeader.Length
	payload := make([]byte, payloadLength)
	_, err = io.ReadFull(r, payload)
	if err != nil {
		return err
	}
	fmt.Printf("  Payload: %v\n", string(payload))

	// TODO QoS0なのでレスポンスなし
	return nil
}
