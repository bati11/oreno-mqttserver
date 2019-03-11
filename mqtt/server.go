package mqtt

import (
	"bufio"
	"fmt"
	"io"
	"net"

	"github.com/bati11/oreno-mqtt/mqtt/handler"
	"github.com/bati11/oreno-mqtt/mqtt/packet"
)

func Run() {
	ln, err := net.Listen("tcp", "localhost:1883")
	if err != nil {
		panic(err)
	}
	fmt.Println("server starts at localhost:1883")

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		err = handle(conn)
		if err != nil {
			panic(err)
		}
	}
}

func handle(conn net.Conn) error {
	defer conn.Close()

	for {
		r := bufio.NewReader(conn)
		mqttReader := packet.NewMQTTReader(r)
		packetType, err := mqttReader.ReadPacketType()
		if err != nil {
			if err == io.EOF {
				// クライアント側から既に切断してる場合
				return nil
			}
			return err
		}
		if packetType == packet.PUBLISH {
			fixedHeader, err := packet.ToPublishFixedHeader(mqttReader)
			if err != nil {
				return err
			}
			fmt.Printf("-----\n%+v\n", fixedHeader)
			err = handler.HandlePublish(fixedHeader, r)
			if err != nil {
				return err
			}
		} else {
			fixedHeader, err := packet.ToFixedHeader(mqttReader)
			if err != nil {
				return err
			}
			fmt.Printf("-----\n%+v\n", fixedHeader)
			switch fixedHeader.PacketType {
			case packet.CONNECT:
				connack, err := handler.HandleConnect(fixedHeader, r)
				if err != nil {
					return err
				}
				_, err = conn.Write(connack.ToBytes())
				if err != nil {
					return err
				}
			case packet.DISCONNECT:
				return nil
			}
		}
	}
}
