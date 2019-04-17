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

		go func() {
			err = handle(conn)
			if err != nil {
				panic(err)
			}
		}()
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
		switch packetType {
		case packet.PUBLISH:
			err = handler.HandlePublish(mqttReader)
			if err != nil {
				return err
			}
		case packet.CONNECT:
			connack, err := handler.HandleConnect(mqttReader)
			if err != nil {
				return err
			}
			_, err = conn.Write(connack.ToBytes())
			if err != nil {
				return err
			}
		case packet.SUBSCRIBE:
			suback, err := handler.HandleSubscribe(mqttReader)
			if err != nil {
				return err
			}
			_, err = conn.Write(suback.ToBytes())
			if err != nil {
				return err
			}
		case packet.PINGREQ:
			pingresp, err := handler.HandlePingreq(mqttReader)
			if err != nil {
				return err
			}
			_, err = conn.Write(pingresp.ToBytes())
			if err != nil {
				return err
			}
		case packet.DISCONNECT:
			return nil
		}
	}
}
