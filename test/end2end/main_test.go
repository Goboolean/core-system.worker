package integration

import (
	"context"
	"testing"
	"time"

	"github.com/Goboolean/core-system.worker/configuration"

	"github.com/Goboolean/core-system.worker/test/container"
	influxutil "github.com/Goboolean/core-system.worker/test/util/influx"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/stretchr/testify/assert"
)

var (
	influxDBOrg   = configuration.InfluxDBOrg
	influxDBToken = configuration.InfluxDBToken

	tradeBucket      = configuration.InfluxDBTradeBucket
	orderBucket      = configuration.InfluxDBOrderEventBucket
	annotationBucket = configuration.InfluxDBAnnotationBucket
)

func TestPing(t *testing.T) {
	rawInfluxClient := influxdb2.NewClient(configuration.InfluxDBURL, influxDBToken)
	t.Cleanup(func() {
		rawInfluxClient.Close()
	})

	ok, err := rawInfluxClient.Ping(context.Background())
	assert.True(t, ok)
	assert.NoError(t, err)
}

func TestCountRecordsInMeasurement(t *testing.T) {
	rawInfluxClient := influxdb2.NewClient(configuration.InfluxDBURL, influxDBToken)
	t.Cleanup(func() {
		rawInfluxClient.Close()
	})

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
	influxC, err := container.InitInfluxContainerWithRandomPort(context.Background(), tradeBucket, orderBucket, annotationBucket)
	if err != nil {
		panic(err)
	}
	configuration.InfluxDBURL = influxC.URL
	m.Run()
	err = influxC.Terminate(context.Background())
	if err != nil {
		panic(err)
	}
}
