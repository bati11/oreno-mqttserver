package handler

import (
	"fmt"

	"github.com/bati11/oreno-mqtt/mqtt/packet"
)

func HandleSubscribe(reader *packet.MQTTReader) (*packet.Suback, error) {
	fmt.Printf("  HandleSubscribe\n")

	subscribe, err := reader.ReadSubscribe()
	if err != nil {
		return nil, err
	}

	fmt.Printf("  %#v\n", subscribe.VariableHeader)
	fmt.Printf("  %+v\n", subscribe.Payload)

	var qoss []uint8
	for range subscribe.Payload.TopicFilterPairs {
		qoss = append(qoss, 0)
	}

	suback := packet.NewSubackForSuccess(subscribe.VariableHeader.PacketIdentifier, qoss)
	return suback, nil
}
