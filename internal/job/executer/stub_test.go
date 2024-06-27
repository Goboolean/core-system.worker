package executer_test

import (
	"testing"

	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/job/executer"
	"github.com/Goboolean/core-system.worker/internal/model"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

func TestStub(t *testing.T) {
	t.Run("stub에 N개의 데이터가 입력됐을 때 N개의 데이터를 출력해야 한다.", func(t *testing.T) {
		//arrange
		num := 100
		inChan := make(job.DataChan, num)

		for i := 0; i < num; i++ {
			inChan <- model.Packet{
				Sequence: int64(i),
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
		close(inChan)
		stub, err := executer.NewStub(&job.UserParams{})
		if err != nil {
			t.Error(err)
			return
		}
		stub.SetInput(inChan)

		//act

		res := make([]model.Packet, 0, num)
		g := errgroup.Group{}

		g.Go(func() error {
			for v := range stub.Output() {
				res = append(res, v)
			}
			return nil
		})

		g.Go(stub.Execute)
		err = g.Wait()

		//assert
		assert.NoError(t, err)
		assert.Len(t, res, num)
	})
}
