package handler

import (
	"bufio"
	"fmt"

	"github.com/bati11/oreno-mqtt/mqtt/packet"
)

func HandlePingreq(fixedHeader packet.FixedHeader, r *bufio.Reader) (packet.Pingresp, error) {
	fmt.Printf("  HandlePingreq\n")
	pingresp := packet.NewPingresp()
	return pingresp, nil
}
