package handler

import (
	"bufio"
	"fmt"

	"github.com/bati11/oreno-mqtt/mqtt/packet"
)

func HandleSubscribe(fixedHeader packet.FixedHeader, r *bufio.Reader) (packet.Suback, error) {
	fmt.Printf("  HandleSubscribe\n")
	variableHeader, err := packet.ToSubscribeVariableHeader(fixedHeader, r)
	if err != nil {
		return packet.Suback{}, err
	}
	fmt.Printf("  %#v\n", variableHeader)

	payload, err := packet.ToSubscribePayload(fixedHeader, variableHeader, r)
	if err != nil {
		return packet.Suback{}, err
	}

	fmt.Printf("  %+v\n", payload)

	var qoss []uint8
	for range payload.TopicFilterPairs {
		qoss = append(qoss, 0)
	}

	suback := packet.NewSubackForSuccess(variableHeader.PacketIdentifier, qoss)
	return suback, nil
}
