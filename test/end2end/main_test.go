package integration

import (
	"context"
	"os"
	"testing"
	"time"

	influxutil "github.com/Goboolean/core-system.worker/test/util/influx"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/stretchr/testify/assert"
)

var rawInfluxClient influxdb2.Client

var (
	influxDBUrl   = os.Getenv("INFLUXDB_URL")
	influxDBOrg   = os.Getenv("INFLUXDB_ORG")
	influxDBToken = os.Getenv("INFLUXDB_TOKEN")

	tradeBucket      = os.Getenv("INFLUXDB_TRADE_BUCKET")
	orderBucket      = os.Getenv("INFLUXDB_ORDER_EVENT_BUCKET")
	annotationBucket = os.Getenv("INFLUXDB_ANNOTATION_BUCKET")
)

func TestPing(t *testing.T) {
	ok, err := rawInfluxClient.Ping(context.Background())
	assert.True(t, ok)
	assert.NoError(t, err)
}

func TestCountRecordsInMeasurement(t *testing.T) {
	bucket := annotationBucket
	measurement := "testMeasurement"
	if err := influxutil.RecreateBucket(rawInfluxClient, influxDBOrg, bucket); err != nil {
		t.Error(err)
		t.FailNow()
	}
	startTime := time.Unix(1720396800, 0)
	num := 390
	writer := rawInfluxClient.WriteAPIBlocking(influxDBOrg, bucket)
	for i := 0; i < num; i++ {
		err := writer.WritePoint(context.Background(),
			write.NewPoint(
				measurement,
				map[string]string{},
				map[string]interface{}{
					"testString": "hello",
					"testNum":    3,
					"testFloat":  3.14,
				},
				startTime.Add(time.Duration(i)*time.Minute),
			))
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
	}

	count, err := influxutil.CountRecordsInMeasurement(rawInfluxClient, influxDBOrg, bucket, measurement)
	assert.Equal(t, num, count)
	assert.NoError(t, err)

}

func TestMain(m *testing.M) {
	rawInfluxClient = influxdb2.NewClient(influxDBUrl, influxDBToken)
	m.Run()
	rawInfluxClient.Close()
}
