package transmitter

import (
	"encoding/json"

	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/model"
	"github.com/Goboolean/core-system.worker/internal/util/chanutil"
	log "github.com/sirupsen/logrus"
)

// Fake logs the data received from the input channel
// without dispatching it.
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
		switch v := in.Data.(type) {
		case *model.TradeCommand:
			log.WithFields(log.Fields{
				"ProportionPercent:": v.ProportionPercent,
				"Action:           ": v.Action.String(),
				"Timestamp:        ": in.Time,
			}).Debug("Order event is dispatched")
		default:
			annotationString, err := json.Marshal(v)
			if err != nil {
				log.Warn(err.Error())
			}
			log.WithField("content", string(annotationString)).Debug(
				"Annotation is dispatched",
			)
		}

	}
	return nil
}

func (f *Fake) SetInput(in job.DataChan) {
	f.in = in
}
