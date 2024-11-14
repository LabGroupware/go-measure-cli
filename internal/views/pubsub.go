package views

import "sync"

type broadcaster[T any] struct {
	subscribers map[chan T]struct{}
	mutex       sync.RWMutex
}

func NewBroadcaster[T any]() *broadcaster[T] {
	return &broadcaster[T]{
		subscribers: make(map[chan T]struct{}),
	}
}

func (b *broadcaster[T]) subscribe() chan T {
	ch := make(chan T)
	b.mutex.Lock()
	b.subscribers[ch] = struct{}{}
	b.mutex.Unlock()
	return ch
}

func (b *broadcaster[T]) unsubscribe(ch chan T) {
	b.mutex.Lock()
	delete(b.subscribers, ch)
	close(ch)
	b.mutex.Unlock()
}

func (b *broadcaster[T]) broadcast(value T) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	for ch := range b.subscribers {
		go func(ch chan T) {
			ch <- value
		}(ch)
	}
}
