package joiner_test

import (
	"sync"
	"testing"
	"time"

	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/job/joiner"
	"github.com/Goboolean/core-system.worker/internal/model"
	"github.com/stretchr/testify/suite"
)

func TestMain(m *testing.M) {
	m.Run()
}

type ByTimeTestSuite struct {
	suite.Suite
}

func (suite *ByTimeTestSuite) TestByTime_ShouldPackDataTogether_WhenDataFromReferenceInputAndFromModelInputHaveSameTime() {
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
			suite.T().Error(err)
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
		suite.NoError(err)
		suite.Equal(exp, res)
	}
}

func TestByTime(t *testing.T) {
	suite.Run(t, new(ByTimeTestSuite))
}
