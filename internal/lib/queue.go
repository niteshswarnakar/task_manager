package lib

type Queue[T any] interface {
	Get() T
	Put(T)
	Close()
	GetChannel() chan T
}

type channelQueue[T any] struct {
	queue chan T
}

func (c channelQueue[T]) Get() T {
	return <-c.queue
}

func (c channelQueue[T]) Put(item T) {
	c.queue <- item
}

func (c channelQueue[T]) Close() {
	close(c.queue)
}

func (c channelQueue[T]) GetChannel() chan T {
	return c.queue
}

func NewQueue[T any]() Queue[T] {
	return channelQueue[T]{
		queue: make(chan T),
	}
}
