package end2end

import (
	"context"
	"testing"
	"time"

	"github.com/Goboolean/core-system.worker/configuration"
	"github.com/Goboolean/core-system.worker/internal/pipeline"
	influxutil "github.com/Goboolean/core-system.worker/test/util/influx"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/stretchr/testify/suite"
)

type WithoutModelTestSuite struct {
	suite.Suite
	rawClient influxdb2.Client
}

func (suite *WithoutModelTestSuite) SetupSuite() {
	suite.rawClient = influxdb2.NewClient(url, influxDBToken)
}

func (suite *WithoutModelTestSuite) TearDownTestSuite() {
	suite.rawClient.Close()
}

func (suite *WithoutModelTestSuite) SetupTest() {
	suite.Require().NoError(influxutil.RecreateBucket(suite.rawClient, influxDBOrg, tradeBucket))
	suite.Require().NoError(influxutil.RecreateBucket(suite.rawClient, influxDBOrg, orderBucket))
	suite.Require().NoError(influxutil.RecreateBucket(suite.rawClient, influxDBOrg, annotationBucket))

}

func (suite *WithoutModelTestSuite) TestWithoutModel_ShouldProcessAllData_WhenVirtualBackTesting() {
	// arrange
	startTime := time.Unix(1720396800, 0)
	num := 350
	writer := suite.rawClient.WriteAPIBlocking(influxDBOrg, tradeBucket)
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

	// act
	config, err := configuration.ImportAppConfigFromFile("./ymls/without_model.test.yml")
	suite.Require().NoError(err)

	p, err := pipeline.Build(*config)
	suite.Require().NoError(err)

	ctx := context.Background()
	err = p.Run(ctx)

	// assert
	suite.NoError(err)

	var count int
	count, err = suite.countMeasurement(annotationBucket, config.TaskID)
	suite.NoError(err)
	suite.Equal(num, count)

	count, err = suite.countMeasurement(annotationBucket, config.TaskID)
	suite.NoError(err)
	suite.Equal(num, count)
}

func (suite *WithoutModelTestSuite) countMeasurement(bucket, measurement string) (int, error) {
	return influxutil.CountRecordsInMeasurement(suite.rawClient, influxDBOrg, bucket, measurement)
}

func TestWithoutNormal(t *testing.T) {
	suite.Run(t, new(NormalTestSuite))
}
