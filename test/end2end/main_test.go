package integration

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Goboolean/core-system.worker/test/container"
	influxutil "github.com/Goboolean/core-system.worker/test/util/influx"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/stretchr/testify/assert"
)

var (
	url           = ""
	influxDBOrg   = os.Getenv("INFLUXDB_ORG")
	influxDBToken = os.Getenv("INFLUXDB_TOKEN")

	tradeBucket      = os.Getenv("INFLUXDB_TRADE_BUCKET")
	orderBucket      = os.Getenv("INFLUXDB_ORDER_EVENT_BUCKET")
	annotationBucket = os.Getenv("INFLUXDB_ANNOTATION_BUCKET")
)

func TestPing(t *testing.T) {
	rawInfluxClient := influxdb2.NewClient(url, influxDBToken)
	t.Cleanup(func() {
		rawInfluxClient.Close()
	})

	ok, err := rawInfluxClient.Ping(context.Background())
	assert.True(t, ok)
	assert.NoError(t, err)
}

func TestCountRecordsInMeasurement(t *testing.T) {
	rawInfluxClient := influxdb2.NewClient(url, influxDBToken)
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
	url = influxC.URL
	// 고언어 테스트가 병렬적으로 실행될 때는 멀티프로세스 환경에서 실행되므로
	// 한 패키지에서 설정한 환경변수는 다른 패키지 테스트에서 영향을 미치지 않는다.
	os.Setenv("INFLUXDB_URL", url)

	m.Run()
	err = influxC.Terminate(context.Background())
	if err != nil {
		panic(err)
	}
}
