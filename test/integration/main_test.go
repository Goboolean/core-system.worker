package integration

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

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

func RecreateBucket(client influxdb2.Client, orgName, bucketName string) error {
	org, err := client.OrganizationsAPI().FindOrganizationByName(context.Background(), orgName)
	if err != nil {
		return err
	}

	bucket, err := client.BucketsAPI().FindBucketByName(context.Background(), bucketName)
	if err != nil {
		return nil
	}

	if err := client.BucketsAPI().DeleteBucket(context.Background(), bucket); err != nil {
		return err
	}
	_, err = client.BucketsAPI().CreateBucketWithName(context.Background(), org, bucketName)

	return err
}

func CountRecordsInMeasurement(client influxdb2.Client, orgName, bucketName, measurement string) (int, error) {

	q, err := client.QueryAPI(orgName).
		Query(context.Background(),
			fmt.Sprintf(
				`from(bucket: "%s")
				|> range(start:0)
				|> filter(fn: (r) => r["_measurement"] == "%s")
				|> count()`, bucketName, measurement))
	if err != nil {
		return 0, err
	}

	num := int64(0)

	// 각 record 별 count에서 최댓값을 찾는다.
	for q.Next() {
		fmt.Println("CountRecordsInMeasurement: records: ", q.Record().Values())
		val := q.Record().ValueByKey("_value").(int64)
		if val > num {
			num = val
		}
	}

	return int(num), nil
}

func TestPing(t *testing.T) {
	ok, err := rawInfluxClient.Ping(context.Background())
	assert.True(t, ok)
	assert.NoError(t, err)
}

func TestCountRecordsInMeasurement(t *testing.T) {
	bucket := annotationBucket
	measurement := "testMeasurement"

	if err := RecreateBucket(rawInfluxClient, influxDBOrg, bucket); err != nil {
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

	count, err := CountRecordsInMeasurement(rawInfluxClient, influxDBOrg, bucket, measurement)
	assert.Equal(t, num, count)
	assert.NoError(t, err)

}

func TestMain(m *testing.M) {
	rawInfluxClient = influxdb2.NewClient(influxDBUrl, influxDBToken)
	m.Run()
	rawInfluxClient.Close()
}
