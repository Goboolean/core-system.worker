package fetcher_test

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/job/fetcher"
	"github.com/Goboolean/core-system.worker/internal/model"
	"github.com/Goboolean/core-system.worker/internal/util"
	"github.com/stretchr/testify/suite"
)

type StubTestSuite struct {
	suite.Suite
}

func (suite *StubTestSuite) TestStub_ShouldOutputRequiredNumOfData() {
	//arrange
	num := 100
	stub, err := fetcher.NewStockStub(&job.UserParams{
		"numOfGeneration":            strconv.FormatInt(int64(num), 10),
		"maxRandomDelayMilliseconds": strconv.FormatInt(100, 10)})
	suite.Require().NoError(err)
	//act

	wg := &sync.WaitGroup{}
	res := make([]model.Packet, 0)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for v := range stub.Output() {
			res = append(res, v)
		}
	}()

	err = stub.Execute()
	if util.IsWaitGroupTimeout(wg, 10*time.Second) {
		suite.T().Errorf("Deadline exceed")
		suite.T().FailNow()
	}

	//assert
	suite.NoError(err, 0)
	suite.Len(res, num)
}

func TestStub(t *testing.T) {
	suite.Run(t, new(StubTestSuite))
}
