package handler

import (
	"fmt"

	"github.com/bati11/oreno-mqtt/mqtt/packet"
)

func HandlePublish(reader *packet.MQTTReader) error {
	fmt.Printf("  HandlePublish\n")
	publish, err := reader.ReadPublish()
	if err != nil {
		return err
	}
	fmt.Printf("  %#v\n", publish.VariableHeader)
	fmt.Printf("  Payload: %v\n", string(publish.Payload))

	// TODO QoS0なのでレスポンスなし
	return nil
}
