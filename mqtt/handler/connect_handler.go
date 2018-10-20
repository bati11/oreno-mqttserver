package handler

import (
	"fmt"

	"github.com/bati11/oreno-mqtt/mqtt/packet"
)

func HandleConnect(fixedHeader packet.FixedHeader, remains []byte) (packet.Connack, error) {
	variableHeader, remains, err := packet.ToConnectVariableHeader(fixedHeader, remains)
	switch err.(type) {
	case *packet.ConnectError:
		fmt.Printf("%#v\n", err)
		return packet.NewConnackForRefusedByUnacceptableProtocolVersion(), nil
	case error:
		return packet.Connack{}, err
	}
	payload, err := packet.ToConnectPayload(remains)
	if err != nil {
		fmt.Printf("%#v\n", err)
		return packet.NewConnackForRefusedByIdentifierRejected(), nil
	}
	fmt.Printf("fixedHeader: %+v\n", fixedHeader)
	fmt.Printf("variableHeader: %+v\n", variableHeader)
	fmt.Printf("payload: %v\n", payload)
	connack := packet.NewConnackForAccepted()
	return connack, nil
}
