package main

import (
	"fmt"

	"github.com/bati11/oreno-mqtt/mqtt/packet"
)

func main() {
	b1 := byte(0x12) // 0001 0 01 0
	b2 := byte(0x00) // 00000000
	in := [2]byte{b1, b2}

	header := packet.ToFixedHeader(in)
	fmt.Println(header)
}
