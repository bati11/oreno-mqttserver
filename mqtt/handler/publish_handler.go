package handler

import (
	"fmt"

	"github.com/bati11/oreno-mqtt/mqtt/packet"
)

func HandlePublish(reader *packet.MQTTReader) (*packet.Publish, error) {
	fmt.Printf("  HandlePublish\n")
	publish, err := reader.ReadPublish()
	if err != nil {
		return nil, err
	}
	fmt.Printf("  %#v\n", publish.VariableHeader)
	fmt.Printf("  Payload: %v\n", string(publish.Payload))

	// TODO QoS0なのでレスポンスなし
	return publish, nil
}
