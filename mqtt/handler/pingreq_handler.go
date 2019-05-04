package handler

import (
	"fmt"

	"github.com/bati11/oreno-mqtt/mqtt/packet"
)

func HandlePingreq(reader *packet.MQTTReader) (*packet.Pingresp, error) {
	fmt.Printf("  HandlePingreq\n")
	pingreq, err := reader.ReadPingreq()
	if err != nil {
		return nil, err
	}
	fmt.Printf("  %+v\n", pingreq)
	pingresp := packet.NewPingresp()
	fmt.Printf("  ---\n")
	fmt.Printf("  %+v\n", pingresp)
	fmt.Printf("  %+v\n", pingresp.ToBytes())
	return &pingresp, nil
}
