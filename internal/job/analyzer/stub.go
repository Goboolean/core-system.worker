package analyzer

import (
	"sync"

	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/model"
	"github.com/Goboolean/core-system.worker/internal/util"
)

type Stub struct {
	Analyzer

	in  chan any
	out chan any

	sn *util.StopNotifier
	wg *sync.WaitGroup
}

func NewStub(parmas *job.UserParams) *Stub {
	instance := &Stub{
		out: make(chan any),
		sn:  util.NewStopNotifier(),
		wg:  &sync.WaitGroup{},
	}

	return instance
}

func (s *Stub) Execute() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer s.sn.NotifyStop()
		defer close(s.out)

		for {
			select {
			case <-s.sn.Done():
				return
			case <-s.in:
				{
					s.out <- model.TransactionDetails{
						Action:            model.Sell,
						ProportionPercent: 100,
					}
				}
			}
		}
	}()
}

func (s *Stub) Close() error {
	s.sn.NotifyStop()
	s.wg.Wait()
	return nil
}

func (s *Stub) SetInput(in chan any) {
	s.in = in
}

func (s *Stub) Output() chan any {
	return s.out
}
