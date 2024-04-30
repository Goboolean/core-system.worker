package fetcher

import (
	"context"
	"sync"
	"time"

	"github.com/Goboolean/core-system.worker/internal/dto"
	"github.com/Goboolean/core-system.worker/internal/infrastructure"
	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/util"
)

type RealtimeStock struct {
	Fetcher

	pastRepo infrastructure.MongoClientStock

	//미리 가져올 데이터의 개수
	prefetchNum int
	timeSlice   string

	in   chan any `type:"none"`
	out  chan any `type:"*StockAggregate"` //Job은 자신의 Output 채널에 대해 소유권을 가진다.
	wg   sync.WaitGroup
	stop *util.StopNotifier
}

func NewFetchRealtimeStockJob(params job.UserParams) *RealtimeStock {
	//여기에 기본값 입력 아웃풋 채널은 job이 소유권을 가져야 한다.
	instance := &RealtimeStock{
		out:  make(chan any),
		stop: util.NewStopNotifier(),
	}

	return instance
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

		//prefetch past stock data
		count := rt.pastRepo.GetCount(ctx)
		duration, _ := time.ParseDuration(rt.timeSlice)

		rt.pastRepo.ForEachDocument(ctx, (count-1)-(rt.prefetchNum), rt.prefetchNum, func(doc infrastructure.StockDocument) {

			rt.out <- &dto.StockAggregate{
				OpenTime:   doc.Timestamp,
				ClosedTime: doc.Timestamp + (duration.Milliseconds() / 1000),
				Open:       doc.Open,
				Closed:     doc.Close,
				High:       doc.High,
				Low:        doc.Low,
				Volume:     float32(doc.Volume),
			}
		})

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

func (rt *RealtimeStock) SetInputChan(input chan any) {
	rt.in = input
}

func (rt *RealtimeStock) OutputChan() chan any {
	return rt.out
}

func (rt *RealtimeStock) Close() error {
	rt.stop.NotifyStop()
	rt.wg.Wait()
	return nil
}
