package mqtt

import (
	"bufio"
	"context"
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
	defer ln.Close()

	fmt.Println("server starts at localhost:1883")

	pub := make(chan *packet.Publish)
	defer close(pub)
	subscriptions := make(chan *Subscription)
	defer close(subscriptions)
	doneSubscriptionToBroker := make(chan *DoneSubscriptionResult)
	defer close(doneSubscriptionToBroker)

	go Broker(pub, subscriptions, doneSubscriptionToBroker)

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		go func() {
			err = handle(conn, pub, subscriptions, doneSubscriptionToBroker)
			if err != nil {
				panic(err)
			}
		}()
	}
}

func handle(conn net.Conn, publishToBroker chan<- *packet.Publish, subscriptionToBroker chan<- *Subscription, doneSubscriptions chan<- *DoneSubscriptionResult) error {
	defer conn.Close()

	var clientID string

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r := bufio.NewReader(conn)
	for {
		mqttReader := packet.NewMQTTReader(r)
		packetType, err := mqttReader.ReadPacketType()
		if err != nil {
			if err == io.EOF {
				fmt.Println("client closed")
				return nil
			}
			return err
		}
		switch packetType {
		case packet.PUBLISH:
			publish, err := handler.HandlePublish(mqttReader)
			if err != nil {
				return err
			}
			publishToBroker <- publish
		case packet.CONNECT:
			connect, connack, err := handler.HandleConnect(mqttReader)
			if err != nil {
				return err
			}
			_, err = conn.Write(connack.ToBytes())
			if err != nil {
				return err
			}
			clientID = connect.Payload.ClientID
		case packet.SUBSCRIBE:
			suback, err := handler.HandleSubscribe(mqttReader)
			if err != nil {
				return err
			}
			_, err = conn.Write(suback.ToBytes())
			if err != nil {
				return err
			}
			subscription, errCh := handleSub(ctx, clientID, conn)
			subscriptionToBroker <- subscription
			go func(ctx context.Context) {
				var result *DoneSubscriptionResult
				select {
				case <-ctx.Done():
					result = NewDoneSubscriptionResult(subscription.clientID, nil)
				case err, ok := <-errCh:
					if !ok {
						return
					}
					result = NewDoneSubscriptionResult(subscription.clientID, err)
				}
				doneSubscriptions <- result
			}(ctx)
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
			fmt.Println("  handle DISCONNECT")
			return nil
		}
	}
	return nil
}

func handleSub(ctx context.Context, clientID string, conn net.Conn) (*Subscription, <-chan error) {
	errCh := make(chan error)
	subscription, pubFromBroker := NewSubscription(clientID)
	go func() {
		defer close(errCh)
		for {
			select {
			case <-ctx.Done():
				return
			case publishedMessage, ok := <-pubFromBroker:
				if !ok {
					return
				}
				bs := publishedMessage.ToBytes()
				_, err := conn.Write(bs)
				if err != nil {
					errCh <- err
				}
			}
		}
	}()
	return subscription, errCh
}
