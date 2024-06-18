package executer

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/Goboolean/core-system.worker/internal/infrastructure/kserve"
	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/model"
	"github.com/Goboolean/core-system.worker/internal/util"
	"github.com/cenkalti/backoff"
)

type Mock struct {
	ModelExecutor

	//user param의 type은 float32
	modelParam1 float32

	batchSize int32
	maxRetry  int32

	kServeClient kserve.Client

	in  job.DataChan `type:""`
	out job.DataChan `type:""` //Job은 자신의 Output 채널에 대해 소유권을 가진다.

	stop *util.StopNotifier
}

func NewMock(kServeClient kserve.Client, params *job.UserParams) (*Mock, error) {
	//여기에 기본값 초기화 아웃풋 채널은 job이 소유권을 가져야 한다.
	instance := &Mock{
		kServeClient: kServeClient,
		maxRetry:     DefaultMaxRetry,
		out:          make(job.DataChan),
		stop:         util.NewStopNotifier(),
	}

	//여기에서 user param 초기화
	if param1, ok := (*params)["param1"]; ok {
		val, err := strconv.ParseFloat(param1, 32)

		if err != nil {
			return nil, fmt.Errorf("create mock model exec job: %w", err)
		}

		instance.modelParam1 = float32(val)
	}

	if param1, ok := (*params)[job.BatchSize]; ok {
		val, err := strconv.ParseInt(param1, 10, 32)

		if err != nil {
			return nil, fmt.Errorf("create mock model exec job: %w", err)
		}

		instance.batchSize = int32(val)
	}

	return instance, nil
}

func (m *Mock) Execute() {

	go func() {
		defer m.stop.NotifyStop()
		defer close(m.out)
		var accumulator = make([]float32, 0)

		for input := range m.in {
			//TODO: 고루틴이 무한정 생성되는 문제 해결
			ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*60*time.Second))
			go func() {
				<-m.stop.Done()
				cancel()
			}()

			data, ok := input.Data.(*model.StockAggregate)

			if !ok {
				panic(fmt.Errorf("model exec job: type mismatch. expected *model.StockAggregate, got %s %w", reflect.TypeOf(input), job.ErrTypeMismatch))
			}

			//데이터를 1차원 텐서 타입으로 변환한다.
			//데이터가 충분히 쌓일 때까지 다음 동작을 실행할 수 없도록 막는다.
			var numOfInput = 4
			accumulator = append(accumulator, data.High, data.Low, data.Open, data.Close)
			if len(accumulator)/numOfInput < int(m.batchSize) {
				continue
			}

			//이를 http client를 이용해 kserve로 보낸다.
			var out []float32

			b := backoff.WithMaxRetries(backoff.WithContext(backoff.NewExponentialBackOff(), ctx), uint64(m.maxRetry))

			if err := backoff.Retry(func() error {
				var err error
				// Shape = [model.StockAggregate에서 사용되는 데이터의 개수 = 7, batch size]
				out, err = m.kServeClient.RequestInference(ctx, []int{numOfInput, int(m.batchSize)}, accumulator)
				return err

			}, b); err != nil {
				panic(fmt.Errorf("model exec job: inference service returns error %w", err))
			}

			accumulator = accumulator[numOfInput:]

			//반환 받은 텐서 타입에서 알맞은 타입으로 가공한다.
			//지금은 모델이 candlestick를 리턴한다고 가정한다.
			//거래량 중요한 데이터가 아니므로 일단 0처리
			m.out <- model.Packet{
				Sequence: input.Sequence,
				Data: &model.StockAggregate{
					OpenTime:   data.ClosedTime,
					ClosedTime: data.ClosedTime + (data.ClosedTime - data.OpenTime),
					High:       out[0],
					Low:        out[1],
					Open:       out[2],
					Close:      out[3],
					Volume:     0.0,
				},
			}

		}
	}()

}

func (m *Mock) SetInput(input job.DataChan) {
	m.in = input
}

func (m *Mock) Output() job.DataChan {
	return m.out
}

func (m *Mock) Cancel() {
	m.stop.NotifyStop()
}
