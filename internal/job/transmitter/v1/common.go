package v1

import (
	"sync"

	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/job/transmitter"
	"github.com/Goboolean/core-system.worker/internal/model"
	"github.com/Goboolean/core-system.worker/internal/util"
)

type Common struct {
	transmitter.Transmitter

	annotationDispatcher transmitter.AnnotationDispatcher
	orderDispatcher      transmitter.OrderEventDispatcher

	in job.DataChan
	sn *util.StopNotifier
	wg *sync.WaitGroup
}

func NewCommon(
	annotationDispatcher transmitter.AnnotationDispatcher,
	orderDispatcher transmitter.OrderEventDispatcher) (*Common, error) {
	return &Common{
		annotationDispatcher: annotationDispatcher,
		orderDispatcher:      orderDispatcher,
		sn:                   util.NewStopNotifier(),
		wg:                   &sync.WaitGroup{},
	}, nil
}

func (b *Common) Execute() {
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		for {
			select {
			case <-b.sn.Done():
				return
			case inPacket, ok := <-b.in:
				if !ok {
					return
				}
				switch v := inPacket.Data.(type) {
				case *model.OrderEvent:
					b.orderDispatcher.Dispatch(v)
				default:
					b.annotationDispatcher.Dispatch(v)
				}
			}
		}
	}()
}

// SetInput sets the input data channel for the transmitter.
func (b *Common) SetInput(in job.DataChan) {
	b.in = in
}

func (b *Common) Close() error {
	b.sn.NotifyStop()
	b.wg.Wait()

	if err := b.orderDispatcher.Close(); err != nil {
		return err
	}

	if err := b.annotationDispatcher.Close(); err != nil {
		return err
	}

	return nil
}
