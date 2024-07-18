package integration

import (
	"context"
	"testing"
	"time"

	"github.com/Goboolean/core-system.worker/configuration"
	"github.com/Goboolean/core-system.worker/internal/pipeline"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/stretchr/testify/suite"
)

type NormalTestSuite struct {
	suite.Suite
	rawClient influxdb2.Client
}

func (suite *NormalTestSuite) SetupSuite() {
	suite.rawClient = influxdb2.NewClient(influxDBUrl, influxDBToken)
}

func (suite *NormalTestSuite) TearDownTestSuite() {
	suite.rawClient.Close()
}

func (suite *NormalTestSuite) SetupTest() {
	suite.Require().NoError(RecreateBucket(rawInfluxClient, influxDBOrg, tradeBucket))
	suite.Require().NoError(RecreateBucket(rawInfluxClient, influxDBOrg, orderBucket))
	suite.Require().NoError(RecreateBucket(rawInfluxClient, influxDBOrg, annotationBucket))

}

func (suite *NormalTestSuite) TestNormal_ShouldProcessAllData_WhenVirtualBackTesting() {
	//arrange
	startTime := time.Unix(1720396800, 0)
	num := 390
	writer := rawInfluxClient.WriteAPIBlocking(influxDBOrg, tradeBucket)
	for i := 0; i < num; i++ {
		err := writer.WritePoint(context.Background(),
			write.NewPoint(
				"stock.aapl.usa.1m",
				map[string]string{},
				map[string]interface{}{
					"open":   float64(i),
					"close":  float64(2.0),
					"high":   float64(3.0),
					"low":    float64(4.0),
					"volume": float64(4),
				},
				startTime.Add(time.Duration(i)*time.Minute),
			))
		suite.Require().NoError(err)
	}

	//act
	config, err := configuration.ImportAppConfigFromFile("./normal.test.yml")
	suite.Require().NoError(err)

	p, err := pipeline.Build(*config)
	suite.Require().NoError(err)

	ctx := context.Background()
	err = p.Run(ctx)

	//assert
	suite.NoError(err)

	var count int
	count, err = suite.countMeasurement(orderBucket, config.TaskID)
	suite.NoError(err)
	suite.Equal(num, count)

	count, err = suite.countMeasurement(annotationBucket, config.TaskID)
	suite.NoError(err)
	suite.Equal(num, count)
}

func (suite *NormalTestSuite) countMeasurement(bucket, measurement string) (int, error) {
	return CountRecordsInMeasurement(suite.rawClient, influxDBOrg, bucket, measurement)
}

func TestNormal(t *testing.T) {
	suite.Run(t, new(NormalTestSuite))
}
