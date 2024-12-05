package utils

import (
	"fmt"
	"sync"
)

type Broadcaster[T any] struct {
	subscribers map[chan T]struct{}
	mutex       sync.RWMutex
}

func NewBroadcaster[T any]() *Broadcaster[T] {
	return &Broadcaster[T]{
		subscribers: make(map[chan T]struct{}),
	}
}

func (b *Broadcaster[T]) Close() {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	for ch := range b.subscribers {
		delete(b.subscribers, ch)
		close(ch)
	}
}

func (b *Broadcaster[T]) Subscribe() <-chan T {
	ch := make(chan T)
	b.mutex.Lock()
	b.subscribers[ch] = struct{}{}
	b.mutex.Unlock()
	return ch
}

func (b *Broadcaster[T]) Unsubscribe(ch chan T) {
	b.mutex.Lock()
	delete(b.subscribers, ch)
	close(ch)
	b.mutex.Unlock()
}

func (b *Broadcaster[T]) Broadcast(value T) <-chan struct{} {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	doneAny := make(chan struct{}, len(b.subscribers))
	done := make(chan struct{})

	go func() {
		defer close(done)
		for range b.subscribers {
			fmt.Println("Broadcasting")
			<-doneAny
		}

	}()

	for ch := range b.subscribers {
		go func(ch chan T) {
			ch <- value
			doneAny <- struct{}{}
		}(ch)
	}

	return done
}
