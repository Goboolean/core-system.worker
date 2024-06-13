package transmitter

import (
	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/model"
	"github.com/Goboolean/core-system.worker/internal/util"
	log "github.com/sirupsen/logrus"
)

// Execute executes the job with the given context.
type Fake struct {
	in job.DataChan

	done *util.StopNotifier
}

func NewFake() (*Fake, error) {
	return &Fake{
		done: util.NewStopNotifier(),
	}, nil

}

func (f *Fake) Execute() {

	go func() {
		defer f.done.NotifyStop()
		for in := range f.in {

			orderEvent := in.Data.(*model.OrderEvent)

			log.WithFields(log.Fields{
				"ProductID:        ": orderEvent.ProductID,
				"ProportionPercent:": orderEvent.Command.ProportionPercent,
				"Action:           ": orderEvent.Command.Action.String(),
				"Timestamp:        ": orderEvent.CreatedAt,
				"Task:             ": orderEvent.Task.String,
			}).Debug("fake event was dispatched")
		}
	}()
}

func (f *Fake) SetInput(in job.DataChan) {
	f.in = in
}

func (f *Fake) Done() chan struct{} {
	return f.done.Done()
}
