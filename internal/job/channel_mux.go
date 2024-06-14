package job

type ChannelMux[T any] struct {
	in  chan T
	out []chan T
}

func (fo *ChannelMux[T]) SetInput(in chan T) {
	fo.in = in
}

func (fo *ChannelMux[T]) Output() chan T {
	ch := make(chan T, 1)
	fo.out = append(fo.out, ch)
	return ch
}

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
