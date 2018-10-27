package mqtt

import (
	"bufio"
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
			r := bufio.NewReader(conn)
			for {
				fixedHeader, err := packet.ToFixedHeader(r)
				if err == io.EOF {
					break
				} else if err != nil {
					// TODO
					panic(err)
				}
				if fixedHeader.RemainingLength > 268435455 {
					// TODO
					panic(err)
				}
				switch fixedHeader.PacketType {
				case packet.CONNECT:
					fmt.Println("----- START CONNECT -----")
					connack, err := handler.HandleConnect(fixedHeader, r)
					if err != nil {
						// TODO
						panic(err)
					}
					conn.Write(connack.ToBytes())
					fmt.Println("----- END CONNECT -----")
				case packet.PUBLISH:
					fmt.Println("----- START PUBLISH -----")
					_, err := handler.HandlePublish(fixedHeader, r)
					if err != nil {
						// TODO
						fmt.Printf("%+v\n", fixedHeader)
						panic(err)
					}
					fmt.Println("----- END PUBLISH -----")
				}
			}
			conn.Close()
		}(conn)
	}
}
