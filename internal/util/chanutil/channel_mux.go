package chanutil

// ChannelMux[T] multiplexes data from a single input channel to multiple output channels concurrently.
// It's useful in scenarios where you need to fan-out data to multiple consumers or workers concurrently.
type ChannelMux[T any] struct {
	in  chan T
	out []chan T
}

// NewChannelMux[T] initializes and returns a new instance of ChannelMux[T].
func NewChannelMux[T any]() *ChannelMux[T] {
	return &ChannelMux[T]{
		out: make([]chan T, 0),
	}
}

// Sets the input channel for the ChannelMux instance.
//
// You MUST SetInput before calling Execute.
func (fo *ChannelMux[T]) SetInput(in chan T) {
	fo.in = in
}

// Output returns the newly created channel, which can be used to receive data from the ChannelMux.
//
// You MUST SetInput before calling Execute.
// To prevent deadlock, all channels returned from Output MUST be consumed.
func (fo *ChannelMux[T]) Output() chan T {
	ch := make(chan T, 1)
	fo.out = append(fo.out, ch)
	return ch
}

// Execute starts a goroutine to continuously read from the input channel (fo.in)
// and distribute data to all output channels (fo.out)
// It closes all output channels (fo.out) when the input channel (fo.in) is closed.
//
// Don't FORGET to call Execute; otherwise, data will not be forwarded from the input channel to the output channels.
func (fo *ChannelMux[T]) Execute() {
	go func() {
		defer func() {
			for _, ch := range fo.out {
				close(ch)
			}
		}()

		for {
			data, ok := <-fo.in
			if !ok {
				return
			}

			for _, ch := range fo.out {
				ch <- data
			}
		}

	}()

}
