package util

type ChannelMux struct {
	in  chan any
	out []chan any

	stop chan struct{}
}

func (fo *ChannelMux) SetInput(in chan any) {
	fo.in = in
}

func (fo *ChannelMux) Output() chan any {
	ch := make(chan any, 1)
	fo.out = append(fo.out, ch)
	return ch
}

func (fo *ChannelMux) Execute() {
	go func() {
		defer func() {
			for _, ch := range fo.out {
				close(ch)
			}
		}()

		select {
		case <-fo.stop:
			return
		case data, ok := <-fo.in:
			if !ok {
				close(fo.stop)
				return
			}

			for _, ch := range fo.out {
				ch <- data
			}
		}

	}()

}
