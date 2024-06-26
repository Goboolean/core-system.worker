package transmitter

import (
	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/model"
	"github.com/Goboolean/core-system.worker/internal/util/chanutil"
	log "github.com/sirupsen/logrus"
)

// Execute executes the job with the given context.
type Fake struct {
	in      job.DataChan
	errChan chan error
}

func NewFake() (*Fake, error) {
	return &Fake{
		errChan: make(chan error),
	}, nil
}

func (f *Fake) Execute() error {
	defer func() { go chanutil.DummyChannelConsumer(f.in) }()
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
	return nil
}

func (f *Fake) SetInput(in job.DataChan) {
	f.in = in
}
