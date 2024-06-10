package analyzer

import (
	"sync"

	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/model"
	"github.com/Goboolean/core-system.worker/internal/util"
)

type Stub struct {
	Analyzer

	in  job.DataChan
	out job.DataChan

	sn *util.StopNotifier
	wg *sync.WaitGroup
}

func NewStub(parmas *job.UserParams) (*Stub, error) {
	instance := &Stub{
		out: make(job.DataChan),
		sn:  util.NewStopNotifier(),
		wg:  &sync.WaitGroup{},
	}

	return instance, nil
}

func (s *Stub) Execute() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer s.sn.NotifyStop()
		defer close(s.out)

		i := 0
		for {

			select {
			case <-s.sn.Done():
				return
			case <-s.in:

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

		}
	}()
}

func (s *Stub) Close() error {
	s.sn.NotifyStop()
	s.wg.Wait()
	return nil
}

func (s *Stub) SetInput(in job.DataChan) {
	s.in = in
}

func (s *Stub) Output() job.DataChan {
	return s.out
}
