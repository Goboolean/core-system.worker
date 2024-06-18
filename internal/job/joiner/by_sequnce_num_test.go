package joiner_test

import (
	"sync"
	"testing"
	"time"

	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/job/joiner"
	"github.com/Goboolean/core-system.worker/internal/model"
	"github.com/Goboolean/core-system.worker/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestJoinBySequnceNum(t *testing.T) {
	t.Run("Sequnce가 같은 두 데이터가 주어졌을 때, 이 두 데이터를 join해서 출력해야 한다.", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			//arrange
			referenceInput := []model.Packet{
				{
					Sequence: 1,
					Data:     1,
				},
				{
					Sequence: 2,
					Data:     1,
				},
				{
					Sequence: 3,
					Data:     1,
				},
			}

			modelInput := []model.Packet{

				{
					Sequence: 2,
					Data:     2,
				},
				{
					Sequence: 3,
					Data:     2,
				},
			}

			exp := []model.Packet{
				{
					Sequence: 2,
					Data: &model.Pair{
						RefData:   1,
						ModelData: 2,
					},
				},
				{
					Sequence: 3,
					Data: &model.Pair{
						RefData:   1,
						ModelData: 2,
					},
				},
			}

			referenceInputChan := make(job.DataChan)
			modelInputChan := make(job.DataChan)

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
			joinJob, err := joiner.NewBySequence(&job.UserParams{})
			if err != nil {
				t.Error(err)
				return
			}

			joinJob.SetRefInput(referenceInputChan)
			joinJob.SetModelInput(modelInputChan)

			//act
			res := make([]model.Packet, 0)
			errsInJob := make([]error, 0)
			joinJob.Execute()

			wg := &sync.WaitGroup{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				for v := range joinJob.Output() {
					res = append(res, v)
				}
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				for v := range joinJob.Error() {
					errsInJob = append(errsInJob, v)
				}
			}()

			if util.IsWaitGroupTimeout(wg, 5*time.Second) {
				t.Error("deadline exceed")
				return
			}

			//assert
			assert.Equal(t, exp, res)
			assert.Len(t, errsInJob, 0)

		}
	})

}
