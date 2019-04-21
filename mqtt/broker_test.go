package mqtt

import (
	"sync"
	"testing"
)

func TestBroker(t *testing.T) {
	publishers := make(chan Publisher)
	subscribers := make(chan Subscriber)

	// "broker" goroutine
	go Broker(publishers, subscribers)

	sub1 := make(chan []byte)
	subscribers <- sub1
	sub2 := make(chan []byte)
	subscribers <- sub2

	// "pub" goroutine
	go func() {
		// 1回目のpublish
		pub1 := make(chan []byte)
		publishers <- pub1
		pub1 <- []byte("hoge")
		close(pub1)

		// 2回目のpublish
		pub2 := make(chan []byte)
		publishers <- pub2
		pub2 <- []byte("fuga")
		close(pub2)
	}()

	wg := sync.WaitGroup{}

	// "sub"
	wg.Add(1)
	go func() {
		defer wg.Done()

		message1_1, ok := <-sub1
		if !ok {
			t.Fatalf("failed read from channel")
		}
		if string(message1_1) != "hoge" {
			t.Fatalf("got %v, want \"hoge\"", string(message1_1))
		}

		message1_2, ok := <-sub1
		if !ok {
			t.Fatalf("failed read from channel")
		}
		if string(message1_2) != "fuga" {
			t.Fatalf("got %v, want \"fuga\"", string(message1_2))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		message2_1, ok := <-sub2
		if !ok {
			t.Fatalf("failed read from channel")
		}
		if string(message2_1) != "hoge" {
			t.Fatalf("got %v, want \"hoge\"", string(message2_1))
		}

		message2_2, ok := <-sub2
		if !ok {
			t.Fatalf("failed read from channel")
		}
		if string(message2_2) != "fuga" {
			t.Fatalf("got %v, want \"fuga\"", string(message2_2))
		}
	}()

	wg.Wait()
	close(publishers)
}
