package analyzer

import (
	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/model"
	"github.com/Goboolean/core-system.worker/internal/util/chanutil"
)

type Stub struct {
	in  job.DataChan
	out job.DataChan
}

func NewStub(parmas *job.UserParams) (*Stub, error) {
	instance := &Stub{
		out: make(job.DataChan),
	}

	return instance, nil
}

func (s *Stub) Execute() error {

	defer close(s.out)
	defer func() {
		go chanutil.DummyChannelConsumer(s.in)
	}()

	i := 0
	for input := range s.in {
		//아무런 동작이 일어나지 않는 값
		s.out <- model.Packet{
			Time: input.Time,
			Data: &model.TradeCommand{
				Action:            model.Sell,
				ProportionPercent: 0,
			},
		}
		i++
	}
	return nil
}

func (s *Stub) SetInput(in job.DataChan) {
	s.in = in
}

func (s *Stub) Output() job.DataChan {
	return s.out
}
