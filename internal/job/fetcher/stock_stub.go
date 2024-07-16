package fetcher

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/model"
	"github.com/Goboolean/core-system.worker/internal/util"
)

// StockStub delivers fake stock data encapsulated in a packet to the output channel.
type StockStub struct {
	numOfGeneration            int
	maxRandomDelayMilliseconds int
	out                        job.DataChan `type:"*StockAggregate"` //Job은 자신의 Output 채널에 대해 소유권을 가진다.
	stop                       *util.StopNotifier
}

// NewStockStub creates new instance of StockStub
//
// Params list:
// "numOfGeneration": The number of data generations
// "maxRandomDelayMilliseconds": Maximum random delay in milliseconds between data generation.
func NewStockStub(parmas *job.UserParams) (*StockStub, error) {
	//여기에 기본값 입력 아웃풋 채널은 job이 소유권을 가져야 한다.

	instance := &StockStub{
		maxRandomDelayMilliseconds: DefaultMaxRandomDelayMilliseconds,
		stop:                       util.NewStopNotifier(),
		out:                        make(job.DataChan),
	}

	if !parmas.IsKeyNilOrEmpty("numOfGeneration") {

		val, err := strconv.ParseInt((*parmas)["numOfGeneration"], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("create past stock fetch job: %w", err)
		}

		instance.numOfGeneration = int(val)

	}

	if !parmas.IsKeyNilOrEmpty("maxRandomDelayMilliseconds") {

		val, err := strconv.ParseInt((*parmas)["maxRandomDelayMilliseconds"], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("create past stock fetch job: %w", err)
		}

		instance.maxRandomDelayMilliseconds = int(val)

	}

	return instance, nil
}

func (ps *StockStub) Execute() error {

	defer close(ps.out)
	start := time.Now()
	for i := 0; i < ps.numOfGeneration; i++ {

		select {
		case <-ps.stop.Done():
			return nil
		default:
			ps.out <- model.Packet{
				Time: start.Add(time.Duration(i) * time.Second),
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

	return nil
}

func (ps *StockStub) Output() job.DataChan {
	return ps.out
}

func (ps *StockStub) NotifyStop() {
	ps.stop.NotifyStop()
}
