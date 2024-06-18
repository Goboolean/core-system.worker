package analyzer

import (
	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/model"
)

type Stub struct {
	in      job.DataChan
	out     job.DataChan
	errChan chan error
}

func NewStub(parmas *job.UserParams) (*Stub, error) {
	instance := &Stub{
		out:     make(job.DataChan),
		errChan: make(chan error),
	}

	return instance, nil
}

func (s *Stub) Execute() {

	go func() {
		defer close(s.errChan)
		defer close(s.out)
		i := 0
		for range s.in {
			//아무런 동작이 일어나지 않는 값
			s.out <- model.Packet{
				Sequence: int64(i),
				Data: &model.TradeCommand{
					Action:            model.Sell,
					ProportionPercent: 0,
				},
			}

			i++
		}
	}()
}

func (s *Stub) SetInput(in job.DataChan) {
	s.in = in
}

func (s *Stub) Output() job.DataChan {
	return s.out
}

func (s *Stub) Error() chan error {
	return s.errChan
}
