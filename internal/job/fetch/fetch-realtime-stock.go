package fetch

import (
	"context"
	"time"

	"github.com/Goboolean/core-system.worker/internal/dto"
	"github.com/Goboolean/core-system.worker/internal/infrastructure"
	"github.com/Goboolean/core-system.worker/internal/job"
)

type RealtimeStock struct {
	Fetcher

	pastRepo infrastructure.MongoClientStock

	//미리 가져올 데이터의 개수
	prefetchNum int
	timeSlice   string

	in  chan any `type:"none"`
	out chan any `type:"*StockAggregate"` //Job은 자신의 Output 채널에 대해 소유권을 가진다.
}

func NewFetchRealtimeStockJob(params job.UserParams) *RealtimeStock {
	//여기에 기본값 입력 아웃풋 채널은 job이 소유권을 가져야 한다.
	instance := &RealtimeStock{
		out: make(chan any),
	}

	return instance
}

func (j *RealtimeStock) Execute(ctx context.Context) {

	go func() {
		defer close(j.out)

		//prefetch past stock data
		count := j.pastRepo.GetCount(ctx)
		duration, _ := time.ParseDuration(j.timeSlice)

		j.pastRepo.ForEachDocument(ctx, (count-1)-(j.prefetchNum), j.prefetchNum, func(doc infrastructure.StockDocument) {

			j.out <- &dto.StockAggregate{
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
			case <-ctx.Done():
				return
				//case <- karfka:

				// 알맞게 변환하기
				// out에다가 던지기
			}
		}
	}()

}

func (j *RealtimeStock) SetInputChan(input chan any) {
	j.in = input
}

func (j *RealtimeStock) OutputChan() chan any {
	return j.out
}
