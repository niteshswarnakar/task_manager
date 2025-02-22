package lib

type Queue[T any] interface {
	Get() T
	Put(T)
}

type ChannelQueue[T any] struct {
	queue chan T
}

func (c ChannelQueue[T]) Get() T {
	return <-c.queue
}

func (c ChannelQueue[T]) Put(item T) {
	c.queue <- item
}

func NewQueue[T any]() Queue[T] {
	return ChannelQueue[T]{
		queue: make(chan T),
	}
}
