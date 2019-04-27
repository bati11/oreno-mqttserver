package mqtt

import "github.com/bati11/oreno-mqtt/mqtt/packet"

type Publisher <-chan *packet.Publish
type Subscription chan<- *packet.Publish

func Broker(publishers <-chan Publisher, subscriptions <-chan Subscription) {
	var ss []Subscription
	for {
		select {
		case sub := <-subscriptions:
			ss = append(ss, sub)
		case pub := <-publishers:
			for message := range pub {

				// FIXME 保持しているSubscriberを走査してメッセージを送る
				// これだと、どれか1つのSubscriberでブロックするとよくない
				// けど、goroutine作ったとしても結局はどこかでブロックして、それが連鎖を引き起こすだけ
				// 多分、バッファをどこかで持つことになる
				for _, sub := range ss {
					sub <- message
				}
			}
		}
	}
}
