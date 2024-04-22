package executer

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"sync"

	"github.com/Goboolean/core-system.worker/internal/dto"
	"github.com/Goboolean/core-system.worker/internal/infrastructure"
	"github.com/Goboolean/core-system.worker/internal/job"
)

type Mock struct {
	ModelExecutor

	//user param의 type은 float32
	modelParam1 float32

	batchSize int32
	maxRetry  int32

	kServeClient infrastructure.KServeClient

	in  chan any `type:`
	out chan any `type:` //Job은 자신의 Output 채널에 대해 소유권을 가진다.

	wg sync.WaitGroup
}

func NewMockModelExecJob(kServeClient infrastructure.KServeClient, params job.UserParams) (*Mock, error) {
	//여기에 기본값 초기화 아웃풋 채널은 job이 소유권을 가져야 한다.
	instance := &Mock{
		maxRetry: 5,
		out:      make(chan any),
	}

	//여기에서 user param 초기화
	if param1, ok := params["param1"]; ok {
		val, err := strconv.ParseFloat(param1, 32)

		if err != nil {
			return nil, err
		}

		instance.modelParam1 = float32(val)
	}

	if param1, ok := params["batchSize"]; ok {
		val, err := strconv.ParseInt(param1, 10, 32)

		if err != nil {
			return nil, err
		}

		instance.batchSize = int32(val)
	}

	return instance, nil
}

func (m *Mock) Execute(ctx context.Context) {

	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		defer close(m.out)
		// Shape = [dto.StockAggregate의 필드 개수 = 7, batch size]
		// 총 데이터 개수 = dto.StockAggregate의 필드 개수(7) * batch size
		acc := make([]float32, m.batchSize*7)

		for {
			select {
			case <-ctx.Done():
				return
			case input, ok := <-m.in:
				if !ok {
					//입력 채널이 닫혔을 때 처리
					return
				}

				data, ok := input.(*dto.StockAggregate)

				if !ok {
					panic(fmt.Errorf("model exec job: type mismatch. expected *dto.StockAggregate, got %s %w", reflect.TypeOf(input), job.TypeMismatchError))
				}

				//데이터를 1차원 텐서 타입으로 변환한다.
				//데이터가 충분히 쌓일 때까지 다음 동작을 실행할 수 없도록 막는다.
				acc = append(acc, data.High, data.Low, data.Open, data.Closed)
				if len(acc) < int(m.batchSize) {
					continue
				}

				//이를 http client를 이용해 kserve로 보낸다.
				var out []float32
				var err error
				for i := 0; i < int(m.maxRetry); i++ {
					out, err = m.kServeClient.RequestInference(ctx, []int{7, int(m.batchSize)}, acc)

					if err == nil {
						break
					}
				}

				if err != nil {
					panic(fmt.Errorf("model exec job: inference service returns error %w", err))
				}

				//반환 받은 텐서 타입에서 알맞은 타입으로 가공한다.
				//지금은 모델이 candlestick를 리턴한다고 가정한다.
				//거래량 중요한 데이터가 아니므로 일단 0처리
				m.out <- &dto.StockAggregate{
					OpenTime:   data.ClosedTime,
					ClosedTime: data.ClosedTime + (data.ClosedTime - data.OpenTime),
					High:       out[0],
					Low:        out[1],
					Open:       out[2],
					Closed:     out[3],
					Volume:     0.0,
				}
			}
		}
	}()

}

func (m *Mock) SetInputChan(input chan any) {
	m.in = input
}

func (m *Mock) OutputChan() chan any {
	return m.out
}

func (m *Mock) Close() {
	m.wg.Wait()
}
