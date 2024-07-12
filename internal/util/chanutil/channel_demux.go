package chanutil

import "sync"

// ChannelDeMux collects data from multiple input channels and forwards it to a single output channel.
// This can be useful when you have multiple sources of data that you want to consolidate into one stream.
type ChannelDeMux[T any] struct {
	in  []chan T
	out chan T
}

// AddInput adds one or more input channels to the ChannelDeMux.
//
// You MUST SetInput before calling Execute.
func NewChannelDeMux[T any]() *ChannelDeMux[T] {
	return &ChannelDeMux[T]{
		in:  make([]chan T, 0),
		out: make(chan T),
	}
}

// Output returns the output channel which consolidates data from all input channels.
//
// You MUST SetInput before calling Execute.
// To prevent deadlock, all channels returned from Output MUST be consumed.
func (dm *ChannelDeMux[T]) AddInput(in ...chan T) {
	dm.in = append(dm.in, in...)
}

func (dm *ChannelDeMux[T]) Output() chan T {
	return dm.out
}

// Execute starts a goroutine to continuously read from the input channels
// and consolidate data to the output channel
// It closes the output channel when the input channels are closed.
//
// Don't FORGET to call Execute; otherwise, data will not be forwarded from the input channels to the output channel.
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
