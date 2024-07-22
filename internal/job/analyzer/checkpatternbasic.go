package analyzer

import (
	"github.com/Goboolean/core-system.worker/internal/algorithm/checkpatternbasic"
	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/model"
	"github.com/Goboolean/core-system.worker/internal/util/chanutil"
)

type CheckPatternBasic struct {
	in  job.DataChan
	out job.DataChan

	model checkpatternbasic.Model
}

func NewCheckPatternBasic(parmas *job.UserParams) (*CheckPatternBasic, error) {
	instance := &CheckPatternBasic{
		out: make(job.DataChan),
	}

	return instance, nil
}

func (s *CheckPatternBasic) Execute() error {
	defer close(s.out)
	defer func() {
		go chanutil.DummyChannelConsumer(s.in)
	}()

	for v := range s.in {
		t := v.Time
		stock := v.Data.(*model.StockAggregate)

		indicator, action := s.model.OnEvent(float64(stock.Close))

		switch action {
		case checkpatternbasic.TradeEventBuy:
			s.out <- model.Packet{
				Time: t,
				Data: &model.TradeCommand{
					Action:            model.Buy,
					ProportionPercent: 0,
				},
			}
		case checkpatternbasic.TradeEventSell:
			s.out <- model.Packet{
				Time: t,
				Data: &model.TradeCommand{
					Action:            model.Sell,
					ProportionPercent: 0,
				},
			}
		}

		s.out <- model.Packet{
			Time: t,
			Data: indicator,
		}
	}
	return nil
}

func (s *CheckPatternBasic) SetInput(in job.DataChan) {
	s.in = in
}

func (s *CheckPatternBasic) Output() job.DataChan {
	return s.out
}
