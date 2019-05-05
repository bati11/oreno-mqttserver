package mqtt

import (
	"sync"
	"testing"

	"github.com/bati11/oreno-mqtt/mqtt/packet"
)

func TestBroker(t *testing.T) {
	pub := make(chan *packet.Publish)
	defer close(pub)
	subscriptions := make(chan *Subscription)
	defer close(subscriptions)
	doneSubscriptionToBroker := make(chan *DoneSubscriptionResult)
	defer close(doneSubscriptionToBroker)

	// "broker" goroutine
	go Broker(pub, subscriptions, doneSubscriptionToBroker)

	sub1, pubFromSubscription1 := NewSubscription("sub1")
	subscriptions <- sub1

	sub2, pubFromSubscription2 := NewSubscription("sub2")
	subscriptions <- sub2

	// "pub" goroutine
	go func() {
		// 1回目のpublish
		pub <- packet.NewPublish("sampleTopic", []byte("hoge"))

		// 2回目のpublish
		pub <- packet.NewPublish("sampleTopic", []byte("fuga"))
	}()

	wg := sync.WaitGroup{}

	// "sub"
	wg.Add(1)
	go func() {
		defer wg.Done()

		message1_1, ok := <-pubFromSubscription1
		if !ok {
			t.Fatalf("failed read from channel")
		}
		if string(message1_1.Payload) != "hoge" {
			t.Fatalf("got %v, want \"hoge\"", string(message1_1.Payload))
		}

		message1_2, ok := <-pubFromSubscription1
		if !ok {
			t.Fatalf("failed read from channel")
		}
		if string(message1_2.Payload) != "fuga" {
			t.Fatalf("got %v, want \"fuga\"", string(message1_2.Payload))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		message2_1, ok := <-pubFromSubscription2
		if !ok {
			t.Fatalf("failed read from channel")
		}
		if string(message2_1.Payload) != "hoge" {
			t.Fatalf("got %v, want \"hoge\"", string(message2_1.Payload))
		}

		message2_2, ok := <-pubFromSubscription2
		if !ok {
			t.Fatalf("failed read from channel")
		}
		if string(message2_2.Payload) != "fuga" {
			t.Fatalf("got %v, want \"fuga\"", string(message2_2.Payload))
		}
	}()

	wg.Wait()
}
