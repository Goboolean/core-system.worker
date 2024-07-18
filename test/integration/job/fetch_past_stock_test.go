package job

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/job/fetcher"
	"github.com/Goboolean/core-system.worker/internal/model"
	influxutil "github.com/Goboolean/core-system.worker/test/util/influx"
	"github.com/Goboolean/fetch-system.IaC/pkg/influx"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PastStockTestSuite struct {
	suite.Suite
	rawClient influxdb2.Client
	query     *influx.DB
	cursor    *fetcher.StockTradeCursor
}

func (suite *PastStockTestSuite) SetupSuite() {
	var err error
	suite.rawClient = influxdb2.NewClient(url, token)

	suite.Require().NoError(err)
}

func (suite *PastStockTestSuite) TearDownSuite() {
	suite.rawClient.Close()
}

func (suite *PastStockTestSuite) SetupTest() {
	var err error

	suite.query, err = influx.NewDB(&influx.Opts{
		URL:             url,
		Token:           token,
		Org:             org,
		TradeBucketName: tradeBucketName,
	})
	suite.Require().NoError(err)

	suite.cursor, err = fetcher.NewStockTradeCursor(suite.query)
	suite.Require().NoError(err)

	err = influxutil.RecreateBucket(suite.rawClient, org, tradeBucketName)
	suite.Require().NoError(err)

}

func (suite *PastStockTestSuite) TestPing() {

	ok, err := suite.rawClient.Ping(context.Background())

	assert.True(suite.T(), ok)
	assert.NoError(suite.T(), err)
}

// "저장된 데이터가 없을 때, 0개의 데이터를 가져와야 한다.
func (suite *PastStockTestSuite) TestPastStock_ShouldNotOutputAnyTrade_WhenDataIsNotStored() {
	//arrange
	start := time.Now()
	fetchJob, err := fetcher.NewPastStock(suite.cursor, &job.UserParams{
		job.ProductID: testStockID,
		job.StartDate: fmt.Sprint(start.Unix()),
		job.EndDate:   fmt.Sprint(start.Add(time.Minute).Unix()),
		job.TimeFrame: "1m",
	})
	suite.Require().NoError(err)

	//act
	out := make([]model.Packet, 0)
	go func() {
		for v := range fetchJob.Output() {
			out = append(out, v)
		}
	}()

	err = fetchJob.Execute()
	//assert
	suite.NoError(err)
	suite.Len(out, 0)
}

// 저장소에 저장된 데이터를 모두 가져와야 한다.
func (suite *PastStockTestSuite) TestPastStock_ShouldFetchAllData_WhenDataIsStored() {
	//arrange
	storeNum := 350
	storeInterval := time.Minute
	start := time.Now().Add(-time.Duration(storeNum) * storeInterval)

	err := suite.insertTestStockData(start, storeInterval, storeNum)
	suite.Require().NoError(err)

	fetchJob, err := fetcher.NewPastStock(suite.cursor, &job.UserParams{
		job.ProductID: testStockID,
		job.StartDate: fmt.Sprint(start.Unix()),
		job.EndDate:   fmt.Sprint(start.Add(time.Duration(storeNum) * storeInterval).Unix()),
		job.TimeFrame: testTimeFrame,
	})

	suite.Require().NoError(err)

	//act
	out := make([]model.Packet, 0)
	go func() {
		for v := range fetchJob.Output() {
			out = append(out, v)
		}
	}()

	err = fetchJob.Execute()
	//assert
	suite.NoError(err)
	suite.Len(out, storeNum)

}

// 존재하지 않는 timeFrame일 때 데이터를 가져와선 안 된다.
func (suite *PastStockTestSuite) TestPastStock_ShouldNotFetchData_WhenTimeFrameDoesNotExist() {
	//arrange

	storeNum := 350
	storeInterval := time.Minute
	start := time.Now().Add(-time.Duration(storeNum) * storeInterval)

	err := suite.insertTestStockData(start, storeInterval, storeNum)
	suite.Require().NoError(err)

	fetchJob, err := fetcher.NewPastStock(suite.cursor, &job.UserParams{
		job.ProductID: testStockID,
		job.StartDate: fmt.Sprint(start.Unix()),
		job.EndDate:   fmt.Sprint(start.Add(time.Duration(storeNum) * storeInterval).Unix()),
		job.TimeFrame: "1h",
	})
	suite.Require().NoError(err)

	//act
	out := make([]model.Packet, 0)
	go func() {
		for v := range fetchJob.Output() {
			out = append(out, v)
		}
	}()

	err = fetchJob.Execute()
	//assert
	suite.NoError(err)
	suite.Len(out, 0)
}

// 존재하지 않는 ProductID일 때 데이터를 가져와선 안 된다.
func (suite *PastStockTestSuite) TestPastStock_shouldNotRetrieveData_whenProductIDDoesNotExist() {
	//arrange
	storeNum := 350
	storeInterval := time.Minute
	start := time.Now().Add(-time.Duration(storeNum) * storeInterval)

	err := suite.insertTestStockData(start, storeInterval, storeNum)
	suite.Require().NoError(err)

	fetchJob, err := fetcher.NewPastStock(suite.cursor, &job.UserParams{
		job.ProductID: "wrongProductID",
		job.StartDate: fmt.Sprint(start.Unix()),
		job.EndDate:   fmt.Sprint(start.Add(time.Duration(storeNum) * storeInterval).Unix()),
		job.TimeFrame: "1h",
	})

	suite.Require().NoError(err)
	//act
	out := make([]model.Packet, 0)
	go func() {
		for v := range fetchJob.Output() {
			out = append(out, v)
		}
	}()

	err = fetchJob.Execute()
	//assert
	suite.NoError(err)
	suite.Len(out, 0)
}

func (suite *PastStockTestSuite) insertTestStockData(start time.Time, interval time.Duration, num int) error {
	writer := suite.rawClient.WriteAPIBlocking(org, tradeBucketName)

	for i := 0; i < num; i++ {
		err := writer.WritePoint(
			context.Background(),
			write.NewPoint(
				fmt.Sprintf("%s.%s", testStockID, testTimeFrame),
				map[string]string{},
				map[string]interface{}{
					"open":   float64(i),
					"close":  float64(2.0),
					"high":   float64(3.0),
					"low":    float64(4.0),
					"volume": float64(4),
				},
				start.Add(time.Duration(i)*interval),
			),
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func TestPastStock(t *testing.T) {
	suite.Run(t, new(PastStockTestSuite))
}
