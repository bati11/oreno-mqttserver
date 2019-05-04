package handler

import (
	"fmt"

	"github.com/bati11/oreno-mqtt/mqtt/packet"
)

func HandleConnect(reader *packet.MQTTReader) (packet.Connack, error) {
	fmt.Printf("HandleConnect\n")
	connect, err := reader.ReadConnect()
	if err != nil {
		if ce, ok := err.(packet.ConnectError); ok {
			return ce.Connack(), nil
		}
		return packet.Connack{}, err
	}

	// TODO variableHeaderとpayloadを使って何かしらの処理
	fmt.Printf("  %#v\n", connect.VariableHeader)
	fmt.Printf("  %#v\n", connect.Payload)

	return packet.NewConnackForAccepted(), nil
}
