package util

import "sync"

type ChannelDeMux[T any] struct {
	in  []chan T
	out chan T
}

func NewChannelDeMux[T any]() *ChannelDeMux[T] {
	return &ChannelDeMux[T]{
		in:  make([]chan T, 0),
		out: make(chan T),
	}
}

func (dm *ChannelDeMux[T]) AddInput(in ...chan T) {
	dm.in = append(dm.in, in...)
}

func (dm *ChannelDeMux[T]) Output() chan T {
	return dm.out
}

func (dm *ChannelDeMux[T]) Execute() {
	wg := &sync.WaitGroup{}
	for _, e := range dm.in {
		wg.Add(1)
		go func(ch chan T) {
			defer wg.Done()

			for v := range ch {
				dm.out <- v
			}

		}(e)
	}

	go func() {
		wg.Wait()
		close(dm.out)
	}()
}
