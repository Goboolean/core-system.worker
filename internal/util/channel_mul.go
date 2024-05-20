package util

import "github.com/Goboolean/core-system.worker/internal/model"

type ChannelMux struct {
	in  chan model.Packet
	out []chan model.Packet

	stop chan struct{}
}

func (fo *ChannelMux) SetInput(in chan model.Packet) {
	fo.in = in
}

func (fo *ChannelMux) Output() chan model.Packet {
	ch := make(chan model.Packet, 1)
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
