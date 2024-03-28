package job

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/Goboolean/core-system.worker/internal/dto"
	"github.com/Goboolean/core-system.worker/internal/infrastructure"
)

var (
	InvalidStockId       = errors.New("fetch: can't parse stockId")
	DocumentTypeMismatch = errors.New("fetch: mongo: document type mismatch")
)

type PastStockFetcher struct {
	Job

	timeSlice           string
	isFetchingFullRange bool
	startTimestamp      int64 // Unix timestamp of start time
	stockId             string
	pastRepo            infrastructure.MongoClientStock

	in  chan any `type:"none"`
	out chan any `type:"*StockAggregate"` //Job은 자신의 Output 채널에 대해 소유권을 가진다.

	err chan error
}

func NewPastStockFetcher(mongo infrastructure.MongoClientStock, parmas *UserParams) (*PastStockFetcher, error) {
	//여기에 기본값 입력 아웃풋 채널은 job이 소유권을 가져야 한다.
	var err error = nil
	instance := &PastStockFetcher{
		timeSlice:           "1m",
		pastRepo:            mongo,
		isFetchingFullRange: true,
		out:                 make(chan any),
		err:                 make(chan error),
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

func (f *PastStockFetcher) Execute(ctx context.Context) chan error {
	errChan := make(chan error)
	go func() {
		defer close(f.out)
		defer close(errChan)
		for {
			select {
			case <-ctx.Done():
				//외부에서 종료 처리가 왔을 때 처리
				return
			default:
				f.pastRepo.SetTarget(f.stockId, f.timeSlice)
				//가져올 데이터의 개수
				var quantity int
				//처음 가져올 데이터의 Index
				var index int
				var err error
				var count int = f.pastRepo.GetCount(ctx)

				if f.isFetchingFullRange {
					index = 0
					quantity = count
				} else {

					index, err = f.pastRepo.FindLatestIndexBy(ctx, f.startTimestamp)
					if err != nil {
						panic(err)
					}

					//1,2,3,4,5
					quantity = count - index
				}

				duration, _ := time.ParseDuration(f.timeSlice)
				err = f.pastRepo.ForEachDocument(ctx, index, quantity, func(doc infrastructure.StockDocument) {
					f.out <- &dto.StockAggregate{
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

				return
			}

		}
	}()
	return errChan
}

func (j *PastStockFetcher) SetInputChan(input chan any) {
	j.in = input
}

func (j *PastStockFetcher) OutputChan() chan any {
	return j.out
}
