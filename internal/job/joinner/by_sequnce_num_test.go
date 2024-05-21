package joinner_test

import (
	"testing"

	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/job/joinner"
	"github.com/Goboolean/core-system.worker/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestJoinBySequnceNum(t *testing.T) {
	t.Run("Sequnce가 같은 두 데이터가 주어졌을 때, 이 두 데이터를 join해서 출력해야 한다.", func(t *testing.T) {
		//arrange
		refrenceInput := []model.Packet{
			{
				Sequnce: 1,
				Data:    1,
			},
			{
				Sequnce: 2,
				Data:    1,
			},
			{
				Sequnce: 3,
				Data:    1,
			},
		}

		modelInput := []model.Packet{

			{
				Sequnce: 2,
				Data:    2,
			},
			{
				Sequnce: 3,
				Data:    2,
			},
		}

		exp := []model.Packet{
			{
				Sequnce: 2,
				Data: &model.Pair{
					RefData:   1,
					ModelData: 2,
				},
			},
			{
				Sequnce: 3,
				Data: &model.Pair{
					RefData:   1,
					ModelData: 2,
				},
			},
		}

		refrenceInputChan := make(job.DataChan)
		modelInputChan := make(job.DataChan)

		go func() {
			defer close(refrenceInputChan)
			for _, e := range refrenceInput {
				refrenceInputChan <- e
			}
		}()

		go func() {
			defer close(modelInputChan)
			for _, e := range modelInput {
				modelInputChan <- e
			}
		}()

		joinner, err := joinner.NewBysequnce(&job.UserParams{})
		if err != nil {
			t.Error(err)
			return
		}

		joinner.SetRefInput(refrenceInputChan)
		joinner.SetModelInput(modelInputChan)
		outputChan := joinner.Output()

		//act
		output := make([]model.Packet, 0)
		joinner.Execute()
		for e := range outputChan {
			output = append(output, e)
		}
		//assert
		assert.Equal(t, exp, output)
	})
}
