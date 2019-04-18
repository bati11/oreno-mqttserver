package mqtt

type Publisher <-chan []byte

func Broker(publishers <-chan Publisher, toSub chan<- []byte) {
	defer close(toSub)
	for fromPub := range publishers {
		for message := range fromPub {
			toSub <- message
		}
	}
}
