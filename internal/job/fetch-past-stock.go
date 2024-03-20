package job

import (
	"context"
	"time"

	"github.com/Goboolean/core-system.worker/internal/dto"
	"github.com/Goboolean/core-system.worker/internal/infrastructure"
)

type PastStockFetcher struct {
	Job

	batchSize int
	timeSlice string

	pastRepo infrastructure.MongoClientStock

	in  chan any `type:"none"`
	out chan any `type:"[]StockAggregate"` //Job은 자신의 Output 채널에 대해 소유권을 가진다.

	err chan error
}

func NewExampleJob(mongo infrastructure.MongoClientStock, parmas UserParams) *PastStockFetcher {
	//여기에 기본값 입력 아웃풋 채널은 job이 소유권을 가져야 한다.
	instance := &PastStockFetcher{
		batchSize: 100,
		timeSlice: "5m",
		pastRepo:  mongo,
		out:       make(chan any),
		err:       make(chan error),
	}

	return instance
}

func (f *PastStockFetcher) Execute(ctx context.Context) chan error {
	errChan := make(chan error)
	go func() {
		defer close(f.out)
		defer close(errChan)
		defer f.pastRepo.Close()

		select {
		case <-ctx.Done():
			//외부에서 종료 처리가 왔을 때 처리
			return
		default:
			count, _ := f.pastRepo.GetCount()
			//몇 개인지 가져온다.

			//제일 과거 데이터부터 batch size만큼 가져온다.
			//만약 남은 데이터가 batch size보다 작으면 종료

			// 궁금증 둘 중에서 어떻게 구현해야 하지?
			//0~99, 100~199, ...?
			//0~100, 1~101, 2~103 ...?
			for i := 0; i < count/f.batchSize; i++ {

				res := make([]dto.StockAggregate, f.batchSize)
				data, err := f.pastRepo.FetchItems(ctx, i*f.batchSize, (i+1)*f.batchSize-1)
				if err != nil {
					errChan <- err
					return
				}

				// TODO: timeslice를 duration으로 변환하는 부분 개발
				for _, element := range data {
					res = append(res, dto.StockAggregate{
						OpenTime:   element.Timestamp,
						ClosedTime: time.Now().Add(5 * 60 * time.Second).Unix(),
						Open:       element.Open,
						Closed:     element.Close,
						Low:        element.Low,
						High:       element.High,
						Volume:     float32(element.Volume),
					})
				}

				f.out <- res
			}

			break
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
