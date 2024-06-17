package fetcher

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Goboolean/core-system.worker/internal/infrastructure/mongo"
	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/model"
	"github.com/Goboolean/core-system.worker/internal/util"
)

type RealtimeStock struct {
	Fetcher

	pastRepo mongo.StockClient

	//미리 가져올 데이터의 개수
	prefetchNum int
	timeSlice   string
	stockID     string

	out  job.DataChan `type:"*StockAggregate"` //Job은 자신의 Output 채널에 대해 소유권을 가진다.
	err  chan error
	wg   sync.WaitGroup
	stop *util.StopNotifier
}

func NewRealtimeStock(mongo mongo.StockClient, params *job.UserParams) (*RealtimeStock, error) {
	//여기에 기본값 입력 아웃풋 채널은 job이 소유권을 가져야 한다.
	instance := &RealtimeStock{
		out:  make(job.DataChan),
		stop: util.NewStopNotifier(),
		err:  make(chan error),
	}

	if !params.IsKeyNilOrEmpty(job.EndDate) {

		val, ok := (*params)[job.EndDate]
		if !ok {
			return nil, fmt.Errorf("create past stock fetch job: %w", ErrInvalidStockID)
		}

		instance.stockID = val
	}

	return instance, nil
}

func (rt *RealtimeStock) Execute() {
	rt.wg.Add(1)
	go func() {
		defer rt.wg.Done()
		defer rt.stop.NotifyStop()
		defer close(rt.out)

		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			<-rt.stop.Done()
			cancel()
		}()

		rt.pastRepo.SetTarget(rt.stockID, rt.timeSlice)

		//prefetch past stock data
		count := rt.pastRepo.GetCount(ctx)
		duration, _ := time.ParseDuration(rt.timeSlice)

		var sequence int64 = 0

		err := rt.pastRepo.ForEachDocument(ctx, (count-1)-(rt.prefetchNum), rt.prefetchNum, func(doc mongo.StockDocument) {

			rt.out <- model.Packet{
				Sequence: sequence,
				Data: &model.StockAggregate{
					OpenTime:   doc.Timestamp,
					ClosedTime: doc.Timestamp + (duration.Milliseconds() / 1000),
					Open:       doc.Open,
					Close:      doc.Close,
					High:       doc.High,
					Low:        doc.Low,
					Volume:     float32(doc.Volume),
				},
			}
			sequence++
		})

		if err != nil {
			panic(err)
		}

		for {
			select {
			case <-rt.stop.Done():
				return
				//case <- karfka:

				// 알맞게 변환하기
				// out에다가 던지기
			}
		}
	}()

}

func (rt *RealtimeStock) Output() job.DataChan {
	return rt.out
}

func (rt *RealtimeStock) Close() error {
	rt.stop.NotifyStop()
	rt.wg.Wait()
	return nil
}

func (rt *RealtimeStock) Error() chan error {
	return rt.err
}
