package executer_test

import (
	"sync"
	"testing"
	"time"

	"github.com/Goboolean/core-system.worker/internal/infrastructure/kserve"
	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/job/executer"
	"github.com/Goboolean/core-system.worker/internal/model"
	"github.com/Goboolean/core-system.worker/internal/util"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestMock(t *testing.T) {
	t.Run("Batch size가 2이상일 때", func(t *testing.T) {
		//arrange
		ctl := gomock.NewController(t)
		m := kserve.NewMockClient(ctl)

		m.EXPECT().RequestInference(gomock.Any(), gomock.Any(), []float32{1, 1, 1, 1, 2, 2, 2, 2}).Return(
			[]float32{1, 2, 3, 4}, nil)
		m.EXPECT().RequestInference(gomock.Any(), gomock.Any(), []float32{2, 2, 2, 2, 3, 3, 3, 3}).Return(
			[]float32{5, 6, 7, 8}, nil)

		input := []*model.StockAggregate{
			{
				OpenTime:   1715329020,
				ClosedTime: 1715329030,
				High:       1,
				Low:        1,
				Open:       1,
				Close:      1,
				Volume:     1,
			}, {
				OpenTime:   1715329030,
				ClosedTime: 1715329040,
				High:       2,
				Low:        2,
				Open:       2,
				Close:      2,
				Volume:     2,
			}, {
				OpenTime:   1715329040,
				ClosedTime: 1715329050,
				High:       3,
				Low:        3,
				Open:       3,
				Close:      3,
				Volume:     3,
			},
		}
		inChan := make(job.DataChan, len(input))
		for i, e := range input {
			inChan <- model.Packet{
				Sequence: int64(i),
				Data:     e,
			}
		}
		close(inChan)

		execute, err := executer.NewMock(m, &job.UserParams{job.BatchSize: "2"})
		if err != nil {
			t.Error(err)
			return
		}
		execute.SetInput(inChan)

		//batch seze가 2이기 때문에 [1,2] [2,3]으로 묶어서 실행이 된다.
		expect := []*model.StockAggregate{
			{
				OpenTime:   1715329040, //미래 예측이므로 out.OpenTime = 두 번째 input.CloseTime
				ClosedTime: 1715329050, //미래 예측이므로 out.CloseTime = 두 번째 input.ClosedTime + (input.ClosedTime - input.OpenTime)
				High:       1,
				Low:        2,
				Open:       3,
				Close:      4,
				Volume:     0,
			}, {
				OpenTime:   1715329050, //미래 예측이므로 out.OpenTime = 세 번째 input.CloseTime
				ClosedTime: 1715329060, //미래 예측이므로 out.CloseTime = 세 번째 input.ClosedTime + (input.ClosedTime - input.OpenTime)
				High:       5,
				Low:        6,
				Open:       7,
				Close:      8,
				Volume:     0,
			},
		}

		//act
		res := []*model.StockAggregate{}
		errsInPipe := make([]error, 0)

		wg := &sync.WaitGroup{}

		wg.Add(1)
		go func() {
			defer wg.Done()
			for v := range execute.Output() {
				res = append(res, v.Data.(*model.StockAggregate))
			}
		}()

		err = execute.Execute()

		if util.IsWaitGroupTimeout(wg, 5*time.Second) {
			t.Error("deadline exceed")
			return
		}
		//assert
		assert.NoError(t, err)
		assert.Equal(t, expect, res)
		assert.Len(t, errsInPipe, 0)
	})
}
