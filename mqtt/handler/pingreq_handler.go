package handler

import (
	"fmt"

	"github.com/bati11/oreno-mqtt/mqtt/packet"
)

func HandlePingreq(fixedHeader packet.FixedHeader) (packet.MQTTPacket, error) {
	if fixedHeader.PacketType != packet.PINGREQ {
		return nil, fmt.Errorf("packet type is not PINGREQ(12). it got is %v", fixedHeader.PacketType)
	}
	result := packet.NewPingresp()
	return &result, nil
}
