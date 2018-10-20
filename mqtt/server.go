package mqtt

import (
	"fmt"
	"net"

	"github.com/bati11/oreno-mqtt/mqtt/handler"
	"github.com/bati11/oreno-mqtt/mqtt/packet"
)

func Run() error {
	ln, err := net.Listen("tcp", "localhost:1883")
	if err != nil {
		return err
	}
	fmt.Println("server starts at localhost:1883")
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		go func() {
			bs := make([]byte, 5)
			n, err := conn.Read(bs)
			if err != nil {
				// TODO
				panic(err)
			}
			if n != 5 {
				// TODO
				panic(err)
			}
			fixedHeader, remains, err := packet.ToFixedHeader(bs[:])
			if err != nil {
				// TODO
				panic(err)
			}
			if fixedHeader.RemainingLength > 268435455 {
				// TODO
				panic(err)
			}
			bs = make([]byte, fixedHeader.RemainingLength)
			_, err = conn.Read(bs)
			if err != nil {
				// TODO
				panic(err)
			}
			for _, b := range bs {
				remains = append(remains, b)
			}
			switch fixedHeader.PacketType {
			case packet.CONNECT:
				connack, err := handler.HandleConnect(fixedHeader, remains)
				if err != nil {
					panic(err)
				}
				conn.Write(connack.ToBytes())
			}
			conn.Close()
		}()
	}
}
