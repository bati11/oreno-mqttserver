package mqtt

import (
	"testing"
)

func TestBroker(t *testing.T) {
	publishers := make(chan Publisher)
	sub := make(chan []byte)

	// "broker" goroutine
	go Broker(publishers, sub)

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

	// "sub"
	message1, ok := <-sub
	if !ok {
		t.Fatalf("failed read from channel")
	}
	if string(message1) != "hoge" {
		t.Fatalf("got %v, want \"hoge\"", string(message1))
	}

	message2, ok := <-sub
	if !ok {
		t.Fatalf("failed read from channel")
	}
	if string(message2) != "fuga" {
		t.Fatalf("got %v, want \"fuga\"", string(message2))
	}

	close(publishers)
}
