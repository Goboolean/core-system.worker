package executer_test

import (
	"time"

	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/job/executer"
	"github.com/Goboolean/core-system.worker/internal/model"
	"github.com/stretchr/testify/suite"
	"golang.org/x/sync/errgroup"
)

type StubTestSuite struct {
	suite.Suite
}

func (suite *StubTestSuite) TestStub_ShouldReturnNDataItems_WhenNItemsAreStubbed() {
	//arrange
	num := 100
	inChan := make(job.DataChan, num)
	start := time.Now()
	for i := 0; i < num; i++ {
		inChan <- model.Packet{
			Time: start.Add(time.Duration(i) * time.Second),
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
	suite.Require().NoError(err)

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
	suite.NoError(err)
	suite.Len(res, num)
}
