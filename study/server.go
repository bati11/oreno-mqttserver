package study

import (
	"bufio"
	"fmt"
	"net"

	"github.com/bati11/oreno-mqtt/study/handler"
	"github.com/bati11/oreno-mqtt/study/packet"
)

func Run() {
	ln, err := net.Listen("tcp", "localhost:1883")
	if err != nil {
		panic(err)
	}
	fmt.Println("server starts at localhost:1883")
	conn, err := ln.Accept()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	r := bufio.NewReader(conn)
	fixedHeader, err := packet.ToFixedHeader(r)
	if err != nil {
		panic(err)
	}

	switch fixedHeader.PacketType {
	case 1:
		connack, err := handler.HandleConnect(fixedHeader, r)
		if err != nil {
			panic(err)
		}
		_, err = conn.Write(connack.ToBytes())
		if err != nil {
			panic(err)
		}
	}
}
