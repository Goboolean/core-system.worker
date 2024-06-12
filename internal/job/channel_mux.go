package job

type ChannelMux struct {
	in  DataChan
	out []DataChan
}

func (fo *ChannelMux) SetInput(in DataChan) {
	fo.in = in
}

func (fo *ChannelMux) Output() DataChan {
	ch := make(DataChan, 1)
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
