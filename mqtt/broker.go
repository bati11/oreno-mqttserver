package mqtt

type Publisher <-chan []byte
type Subscriber chan<- []byte

func Broker(publishers <-chan Publisher, subscribers <-chan Subscriber) {
	var ss []Subscriber
	for {
		select {
		case sub := <-subscribers:
			ss = append(ss, sub)
		case pub := <-publishers:
			for message := range pub {
				for _, sub := range ss {
					sub <- message
				}
			}
		}
	}
}
