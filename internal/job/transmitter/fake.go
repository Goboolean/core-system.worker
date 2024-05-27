package transmitter

import (
	"fmt"
	"sync"

	"github.com/Goboolean/core-system.worker/internal/model"
	"github.com/Goboolean/core-system.worker/internal/util"
)

// Execute executes the job with the given context.
type Fake struct {
	Transmitter

	in chan any

	wg *sync.WaitGroup
	sn *util.StopNotifier
}

func NewFake() (*Fake, error) {
	return &Fake{
		wg: &sync.WaitGroup{},
		sn: util.NewStopNotifier(),
	}, nil

}

func (f *Fake) Execute() {
	f.wg.Add(1)

	go func() {
		defer f.wg.Done()
		select {
		case <-f.sn.Done():
		case in, ok := <-f.in:
			if !ok {
				return
			}

			orderEvent := in.(*model.OrderEvent)

			fmt.Println("ProductID:        ", orderEvent.ProductID)
			fmt.Println("ProportionPercent:", orderEvent.Trade.ProportionPercent)
			fmt.Println("Action:           ", orderEvent.Trade.Action.String())
			fmt.Println("Timestamp:        ", orderEvent.Timestamp)
			fmt.Println("task:             ", orderEvent.Task.String())
		}
	}()
}

func (f *Fake) Close() error {
	f.sn.NotifyStop()
	f.wg.Wait()
	return nil
}

func (f *Fake) SetInput(in chan any) {
	f.in = in
}
