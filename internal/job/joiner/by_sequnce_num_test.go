package joiner_test

import (
	"sync"
	"testing"
	"time"

	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/job/joiner"
	"github.com/Goboolean/core-system.worker/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestJoinBySequnceNum(t *testing.T) {
	t.Run("Sequnce가 같은 두 데이터가 주어졌을 때, 이 두 데이터를 join해서 출력해야 한다.", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			start := time.Now()
			//arrange
			referenceInput := []model.Packet{
				{
					Time: start,
					Data: 1,
				},
				{
					Time: start.Add(time.Duration(1) * time.Second),
					Data: 1,
				},
				{
					Time: start.Add(time.Duration(2) * time.Second),
					Data: 1,
				},
			}

			modelInput := []model.Packet{

				{
					Time: start.Add(time.Duration(1) * time.Second),
					Data: 2,
				},
				{
					Time: start.Add(time.Duration(2) * time.Second),
					Data: 2,
				},
			}

			exp := []model.Packet{
				{
					Time: start.Add(time.Duration(1) * time.Second),
					Data: &model.Pair{
						RefData:   1,
						ModelData: 2,
					},
				},
				{
					Time: start.Add(time.Duration(2) * time.Second),
					Data: &model.Pair{
						RefData:   1,
						ModelData: 2,
					},
				},
			}

			referenceInputChan := make(job.DataChan)
			modelInputChan := make(job.DataChan)
			// Insert data into the input channel using a goroutines to insert randomly
			go func() {
				defer close(referenceInputChan)
				for _, e := range referenceInput {
					referenceInputChan <- e
				}
			}()
			go func() {
				defer close(modelInputChan)
				for _, e := range modelInput {
					modelInputChan <- e
				}
			}()

			joinJob, err := joiner.NewByTime(&job.UserParams{})
			if err != nil {
				t.Error(err)
				return
			}

			joinJob.SetRefInput(referenceInputChan)
			joinJob.SetModelInput(modelInputChan)

			//act
			res := make([]model.Packet, 0)

			wg := &sync.WaitGroup{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				for v := range joinJob.Output() {
					res = append(res, v)
				}
			}()

			err = joinJob.Execute()
			wg.Wait()

			//assert
			assert.NoError(t, err)
			assert.Equal(t, exp, res)
		}
	})

}
