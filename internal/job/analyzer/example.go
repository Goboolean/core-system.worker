package analyzer

import (
	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/model"
	"github.com/Goboolean/core-system.worker/internal/util/chanutil"
)

type Example struct {
	in  job.DataChan
	out job.DataChan
}

func NewExample(parmas *job.UserParams) (*Example, error) {
	instance := &Example{
		out: make(job.DataChan),
	}

	return instance, nil
}

func (s *Example) Execute() error {

	defer close(s.out)
	defer func() {
		go chanutil.DummyChannelConsumer(s.in)
	}()

	for v := range s.in {
		t := v.Time
		//stock := v.Data.(*model.StockAggregate)
		//여기에 연산 로직 구현

		s.out <- model.Packet{
			Time: t,
			Data: &model.TradeCommand{
				Action:            model.Sell,
				ProportionPercent: 0,
			},
		}
	}
	return nil
}

func (s *Example) SetInput(in job.DataChan) {
	s.in = in
}

func (s *Example) Output() job.DataChan {
	return s.out
}
