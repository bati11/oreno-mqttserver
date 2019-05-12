package handler

import (
	"fmt"

	"github.com/bati11/oreno-mqtt/mqtt/packet"
)

func HandleConnect(reader *packet.MQTTReader) (*packet.Connect, *packet.Connack, error) {
	fmt.Printf("HandleConnect\n")
	connect, err := reader.ReadConnect()
	if err != nil {
		if ce, ok := err.(packet.ConnectError); ok {
			connack := ce.Connack()
			return connect, &connack, nil
		}
		return connect, &packet.Connack{}, err
	}

	// TODO variableHeaderとpayloadを使って何かしらの処理
	fmt.Printf("  %#v\n", connect.VariableHeader)
	fmt.Printf("  %#v\n", connect.Payload)

	connack := packet.NewConnackForAccepted()
	return connect, &connack, nil
}
