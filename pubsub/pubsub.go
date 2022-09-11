package pubsub

import (
	"github.com/pixlcrashr/squadmortar.xyz-sessions-server/eventhandler"
	"sync"
)

type Publisher[V any] interface {
	Publish(value V)
}

type Subscriber[V any] interface {
	Subscribe() Subscription[V]
}

type Subject[V any] interface {
	Publisher[V]
	Subscriber[V]
}

type Subscription[V any] interface {
	Chan() <-chan V
	Unsubscribe()
}

type subject[V any] struct {
	subscribers []chan V

	mtx sync.RWMutex
}

type subscription[V any] struct {
	ch <-chan V

	unsubscribed bool

	unsubscribeEventHandler eventhandler.EventHandler[*subscription[V], struct{}]
}

func (s *subscription[V]) Chan() <-chan V {
	return s.ch
}

func (s *subscription[V]) Unsubscribe() {
	if s.unsubscribed {
		return
	}

	s.unsubscribed = true
	s.unsubscribeEventHandler.Invoke(s, struct{}{})
}

func NewSubject[V any]() Subject[V] {
	return &subject[V]{
		subscribers: make([]chan V, 0),
		mtx:         sync.RWMutex{},
	}
}

func (s *subject[V]) Publish(value V) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	for _, sub := range s.subscribers {
		sub <- value
	}
}

func (s *subject[V]) Subscribe() Subscription[V] {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	ch := make(chan V, 10)

	s.subscribers = append(s.subscribers, ch)

	sub := &subscription[V]{
		unsubscribed:            false,
		ch:                      ch,
		unsubscribeEventHandler: eventhandler.New[*subscription[V], struct{}](),
	}

	sub.unsubscribeEventHandler.Add(s.unsubscribeEvent)

	return sub
}

func (s *subject[V]) unsubscribeEvent(sub *subscription[V], _ struct{}) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	idx := s.findChIndex(sub.ch, s.subscribers)
	if idx == -1 {
		return
	}

	ch := s.subscribers[idx]
	close(ch)

	s.subscribers[idx] = s.subscribers[len(s.subscribers)-1]
	s.subscribers = s.subscribers[:len(s.subscribers)-1]

	sub.unsubscribeEventHandler.Remove(s.unsubscribeEvent)
}

func (s *subject[V]) findChIndex(ch <-chan V, chs []chan V) int {
	idx := 0

	for _, v := range chs {
		if ch == v {
			return idx
		}

		idx++
	}

	idx = -1
	return idx
}
