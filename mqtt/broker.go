package mqtt

import (
	"fmt"
	"sync"

	"github.com/bati11/oreno-mqtt/mqtt/packet"
)

type Subscription struct {
	clientID string
	pubToSub chan<- *packet.Publish
}

func NewSubscription(clientID string) (*Subscription, <-chan *packet.Publish) {
	pub := make(chan *packet.Publish)
	s := &Subscription{
		clientID: clientID,
		pubToSub: pub,
	}
	return s, pub
}

type DoneSubscriptionResult struct {
	clientID string
	err      error
}

func NewDoneSubscriptionResult(clientID string, err error) *DoneSubscriptionResult {
	return &DoneSubscriptionResult{clientID, err}
}

type subscriptionMap struct {
	syncMap sync.Map
}

func newSubscriptionMap() *subscriptionMap {
	return &subscriptionMap{}
}

func (m *subscriptionMap) get(clientID string) *Subscription {
	s, ok := m.syncMap.Load(clientID)
	if !ok {
		return nil
	}
	return s.(*Subscription)
}

func (m *subscriptionMap) put(clientID string, s *Subscription) {
	m.syncMap.Store(clientID, s)
}

func (m *subscriptionMap) delete(clientID string) {
	m.syncMap.Delete(clientID)
}

func (m *subscriptionMap) apply(f func(s *Subscription)) {
	m.syncMap.Range(func(k, v interface{}) bool {
		s := v.(*Subscription)
		f(s)
		return true
	})
}

func Broker(fromPub <-chan *packet.Publish, subscriptions <-chan *Subscription, doneSubscriptions <-chan *DoneSubscriptionResult) {
	sMap := newSubscriptionMap()
	for {
		select {
		case sub := <-subscriptions:
			sMap.put(sub.clientID, sub)
		case message := <-fromPub:
			// FIXME 保持しているSubscriberを走査してメッセージを送る
			// これだと、どれか1つのSubscriberでブロックするとよくない
			// けど、goroutine作ったとしても結局はどこかでブロックして、それが連鎖を引き起こすだけ
			// 多分、バッファをどこかで持つことになる
			sMap.apply(func(sub *Subscription) {
				sub.pubToSub <- message
			})
		case done := <-doneSubscriptions:
			fmt.Printf("close subscription: %v\n", done.clientID)
			if done.err != nil {
				fmt.Println(done.err)
			}
			sub := sMap.get(done.clientID)
			if sub != nil {
				close(sub.pubToSub)
				sMap.delete(done.clientID)
			}
		}
	}
}
