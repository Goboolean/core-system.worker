package executer_test

import (
	"testing"

	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/job/executer"
	"github.com/Goboolean/core-system.worker/internal/model"
	"github.com/stretchr/testify/assert"
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
		stub.SetInput(inChan)
		//act
		stub.Execute()
		out := make([]model.Packet, 0, num)
		outchan := stub.Output()
		for e := range outchan {
			out = append(out, e)
		}

		//assert
		assert.NoError(t, err)
		assert.Len(t, out, num)
	})
}
