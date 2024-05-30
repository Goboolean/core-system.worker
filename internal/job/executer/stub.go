package executer

import (
	"sync"

	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/model"
	"github.com/Goboolean/core-system.worker/internal/util"
)

type Stub struct {
	ModelExecutor

	in  job.DataChan `type:""`
	out job.DataChan `type:""` //Job은 자신의 Output 채널에 대해 소유권을 가진다.

	wg sync.WaitGroup

	stop *util.StopNotifier
}

func NewStub(params *job.UserParams) (*Stub, error) {
	//여기에 기본값 초기화 아웃풋 채널은 job이 소유권을 가져야 한다.
	instance := &Stub{
		out:  make(job.DataChan),
		stop: util.NewStopNotifier(),
	}

	return instance, nil
}

func (m *Stub) Execute() {

	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		defer m.stop.NotifyStop()
		defer close(m.out)

		for {
			select {
			case <-m.stop.Done():
				return
			case input, ok := <-m.in:
				if !ok {
					//입력 채널이 닫혔을 때 처리
					return
				}

				data := input.Data.(*model.StockAggregate)

				m.out <- model.Packet{
					Sequnce: input.Sequnce,
					Data: &model.StockAggregate{
						OpenTime:   data.ClosedTime,
						ClosedTime: data.ClosedTime + (data.ClosedTime - data.OpenTime),
						High:       data.High,
						Low:        data.Low,
						Open:       data.Open,
						Close:      data.Close,
						Volume:     0.0,
					},
				}
			}
		}
	}()

}

func (m *Stub) SetInput(input job.DataChan) {
	m.in = input
}

func (m *Stub) Output() job.DataChan {
	return m.out
}

func (m *Stub) Close() error {
	m.stop.NotifyStop()
	m.wg.Wait()
	return nil
}
