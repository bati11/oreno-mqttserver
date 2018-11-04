package handler

import (
	"bufio"
	"fmt"
	"io"

	"github.com/bati11/oreno-mqtt/mqtt/packet"
)

func HandleSubscribe(fixedHeader packet.FixedHeader, r *bufio.Reader) (packet.MQTTPacket, error) {
	if fixedHeader.PacketType != packet.SUBSCRIBE {
		return nil, fmt.Errorf("packet type is not SUBSCRIBE(8). it got is %v", fixedHeader.PacketType)
	}
	variableHeaderLength := 2
	bs := make([]byte, variableHeaderLength)
	_, err := io.ReadFull(r, bs)
	if err != nil {
		return nil, err
	}
	variableHeader, err := packet.ToSubscribeVariableHeader(bs)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%+v\n", variableHeader)

	payloadLength := fixedHeader.RemainingLength - uint(variableHeaderLength)
	lReader := io.LimitReader(r, int64(payloadLength))
	subscribePayload, err := packet.ToSubscribePayload(lReader)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%+v\n", subscribePayload)

	result := packet.NewSuback(variableHeader, subscribePayload)
	return &result, nil
}
