package handler

import (
	"bufio"
	"fmt"

	"github.com/bati11/oreno-mqtt/study/packet"
)

// CONNECTパケットの可変ヘッダーのバイト数
var variableHeaderLength = 10

func HandleConnect(fixedHeader packet.FixedHeader, r *bufio.Reader) (packet.Connack, error) {
	variableHeader, err := packet.ToConnectVariableHeader(fixedHeader, r)
	if err != nil {
		if ce, ok := err.(packet.ConnectError); ok {
			return ce.Connack(), nil
		}
		return packet.Connack{}, err
	}

	payload, err := packet.ToConnectPayload(r)
	if err != nil {
		if ce, ok := err.(packet.ConnectError); ok {
			return ce.Connack(), nil
		}
		return packet.Connack{}, err
	}

	// TODO variableHeaderとpayloadを使って何かしらの処理
	fmt.Printf("%+v\n", variableHeader)
	fmt.Printf("%+v\n", payload)

	return packet.NewConnackForAccepted(), nil
}
