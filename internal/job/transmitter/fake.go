package transmitter

import (
	"sync"

	"github.com/Goboolean/core-system.worker/internal/model"
	"github.com/Goboolean/core-system.worker/internal/util"
	log "github.com/sirupsen/logrus"
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

			log.WithFields(log.Fields{
				"ProductID:        ": orderEvent.ProductID,
				"ProportionPercent:": orderEvent.Transaction.ProportionPercent,
				"Action:           ": orderEvent.Transaction.Action.String(),
				"Timestamp:        ": orderEvent.Timestamp,
				"Task:             ": orderEvent.Task.String,
			}).Debug("fake event was dispatched")
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
