package fetcher

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/Goboolean/core-system.worker/internal/dto"
	"github.com/Goboolean/core-system.worker/internal/infrastructure"
	"github.com/Goboolean/core-system.worker/internal/job"
)

var (
	InvalidStockId       = errors.New("fetch: can't parse stockId")
	DocumentTypeMismatch = errors.New("fetch: mongo: document type mismatch")
)

type PastStock struct {
	Fetcher

	timeSlice           string
	isFetchingFullRange bool
	startTimestamp      int64 // Unix timestamp of start time
	stockId             string
	pastRepo            infrastructure.MongoClientStock

	in  chan any `type:"none"`
	out chan any `type:"*StockAggregate"` //Job은 자신의 Output 채널에 대해 소유권을 가진다.

}

func NewPastStockFetcher(mongo infrastructure.MongoClientStock, parmas *job.UserParams) (*PastStock, error) {
	//여기에 기본값 입력 아웃풋 채널은 job이 소유권을 가져야 한다.
	var err error = nil
	instance := &PastStock{
		timeSlice:           "1m",
		pastRepo:            mongo,
		isFetchingFullRange: true,
		out:                 make(chan any),
	}

	if !parmas.IsKeyNullOrEmpty("stockId") {

		val, ok := (*parmas)["stockId"]
		if !ok {
			return nil, InvalidStockId
		}

		instance.stockId = val
	}

	if !parmas.IsKeyNullOrEmpty("startDate") {
		instance.isFetchingFullRange = false

		val, err := strconv.ParseInt((*parmas)["startDate"], 10, 64)
		if err != nil {
			return nil, err
		}

		instance.startTimestamp = val

	}

	return instance, err
}

func (ps *PastStock) Execute(ctx context.Context) {
	go func() {
		defer close(ps.out)

		ps.pastRepo.SetTarget(ps.stockId, ps.timeSlice)
		//가져올 데이터의 개수
		var quantity int
		//처음 가져올 데이터의 Index
		var index int
		var err error
		var count int = ps.pastRepo.GetCount(ctx)

		if ps.isFetchingFullRange {
			index = 0
			quantity = count
		} else {

			index, err = ps.pastRepo.FindLatestIndexBy(ctx, ps.startTimestamp)
			if err != nil {
				panic(err)
			}

			//1,2,3,4,5
			quantity = count - index
		}

		duration, _ := time.ParseDuration(ps.timeSlice)
		err = ps.pastRepo.ForEachDocument(ctx, index, quantity, func(doc infrastructure.StockDocument) {
			ps.out <- &dto.StockAggregate{
				OpenTime:   doc.Timestamp,
				ClosedTime: doc.Timestamp + (duration.Milliseconds() / 1000),
				Open:       doc.Open,
				Closed:     doc.Close,
				High:       doc.High,
				Low:        doc.Low,
				Volume:     float32(doc.Volume),
			}
		})
		if err != nil {
			panic(err)
		}

	}()
}

func (ps *PastStock) Output() chan any {
	return ps.out
}