package mqtt

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/bati11/oreno-mqtt/mqtt/handler"
	"github.com/bati11/oreno-mqtt/mqtt/packet"
	"github.com/gorilla/websocket"
)

func Run(withWebSocket bool) {
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

	if withWebSocket {
		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			Subprotocols:    []string{"mqtt"},
			CheckOrigin: func(r *http.Request) bool {
				// FIXME
				return true
			},
		}
		go func() {
			fmt.Println("websocket server starts at localhost:9090")
			http.HandleFunc("/",
				func(w http.ResponseWriter, r *http.Request) {
					conn, err := upgrader.Upgrade(w, r, nil)
					if err != nil {
						fmt.Println(err)
						return
					}
					mqttConn := packet.NewMQTTConnWithWebSocket(conn)
					handle(mqttConn, pub, subscriptions, doneSubscriptionToBroker)
				})
			err = http.ListenAndServe(":9090", nil)
			if err != nil {
				panic("ListenAndServe: " + err.Error())
			}
		}()
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		mqttConn := packet.NewMQTTConn(conn)

		go func() {
			err = handle(mqttConn, pub, subscriptions, doneSubscriptionToBroker)
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
		}()
	}
}

func handle(conn *packet.MQTTConn, publishToBroker chan<- *packet.Publish, subscriptionToBroker chan<- *Subscription, doneSubscriptions chan<- *DoneSubscriptionResult) error {
	defer conn.Close()

	var clientID string

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		mqttReader := packet.NewMQTTReader(conn)
		packetType, err := mqttReader.ReadPacketType()
		if err != nil {
			if err != io.EOF {
				fmt.Printf("failed to mqttReader.ReadPacketType. err: %v\n", err)
				return err
			}
			fmt.Println("client EOF")
			return nil
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
			err = conn.Write(connack.ToBytes())
			if err != nil {
				return err
			}
			clientID = connect.Payload.ClientID
		case packet.SUBSCRIBE:
			suback, err := handler.HandleSubscribe(mqttReader)
			if err != nil {
				return err
			}
			err = conn.Write(suback.ToBytes())
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
			err = conn.Write(pingresp.ToBytes())
			if err != nil {
				return err
			}
		case packet.DISCONNECT:
			fmt.Println("  handle DISCONNECT")
		}
	}
}

func handleSub(ctx context.Context, clientID string, conn *packet.MQTTConn) (*Subscription, <-chan error) {
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
				err := conn.Write(bs)
				if err != nil {
					errCh <- err
				}
			}
		}
	}()
	return subscription, errCh
}
