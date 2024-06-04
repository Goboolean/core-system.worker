package fetcher

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/model"
	"github.com/Goboolean/core-system.worker/internal/util"
)

type StockStub struct {
	Fetcher

	numOfGeneration            int
	maxRandomDelayMilliseconds int

	out job.DataChan `type:"*StockAggregate"` //Job은 자신의 Output 채널에 대해 소유권을 가진다.

	wg   sync.WaitGroup
	stop *util.StopNotifier
}

func NewStockStub(parmas *job.UserParams) (*StockStub, error) {
	//여기에 기본값 입력 아웃풋 채널은 job이 소유권을 가져야 한다.

	instance := &StockStub{
		maxRandomDelayMilliseconds: DefaultMaxRandomDelayMilliseconds,
		stop:                       util.NewStopNotifier(),
		out:                        make(job.DataChan),
	}

	if !parmas.IsKeyNullOrEmpty("numOfGeneration") {

		val, err := strconv.ParseInt((*parmas)["numOfGeneration"], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("create past stock fetch job: %w", err)
		}

		instance.numOfGeneration = int(val)

	}

	if !parmas.IsKeyNullOrEmpty("maxRandomDelayMilliseconds") {

		val, err := strconv.ParseInt((*parmas)["maxRandomDelayMilliseconds"], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("create past stock fetch job: %w", err)
		}

		instance.maxRandomDelayMilliseconds = int(val)

	}

	return instance, nil
}

func (ps *StockStub) Execute() {
	ps.wg.Add(1)
	go func() {
		defer ps.wg.Done()
		defer ps.stop.NotifyStop()
		defer close(ps.out)

		for i := 0; i < ps.numOfGeneration; i++ {

			select {
			case <-ps.stop.Done():
				return
			default:
				ps.out <- model.Packet{
					Sequence: int64(i),
					Data: &model.StockAggregate{
						OpenTime:   1716775499,
						ClosedTime: 1716775499,
						Open:       1.0,
						Close:      2.0,
						High:       3.0,
						Low:        4.0,
						Volume:     5.0,
					},
				}
			}

			if ps.maxRandomDelayMilliseconds > 0 {
				time.Sleep(time.Duration(rand.Intn(ps.maxRandomDelayMilliseconds)) * time.Millisecond)
			}
		}

	}()
}

func (ps *StockStub) Output() job.DataChan {
	return ps.out
}

func (ps *StockStub) Close() error {
	ps.stop.NotifyStop()
	ps.wg.Wait()
	return nil
}
