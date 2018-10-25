package mqtt

import (
	"fmt"
	"io"
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
		go func(conn net.Conn) {
			for {
				bs := make([]byte, 5)
				n, err := conn.Read(bs)
				if err == io.EOF {
					break
				} else if err != nil {
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
					fmt.Println("----- START CONNECT -----")
					connack, err := handler.HandleConnect(fixedHeader, remains)
					if err != nil {
						// TODO
						panic(err)
					}
					conn.Write(connack.ToBytes())
					fmt.Println("----- END CONNECT -----")
				case packet.PUBLISH:
					fmt.Println("----- START PUBLISH -----")
					_, err := handler.HandlePublish(fixedHeader, remains)
					if err != nil {
						// TODO
						panic(err)
					}
					fmt.Println("----- END PUBLISH -----")
				}
			}
			conn.Close()
		}(conn)
	}
}
