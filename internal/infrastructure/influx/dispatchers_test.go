package influx_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Goboolean/core-system.worker/internal/infrastructure/influx"
	"github.com/Goboolean/core-system.worker/internal/model"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/stretchr/testify/assert"
)

var (
	url    = os.Getenv("INFLUXDB_URL")
	bucket = os.Getenv("INFLUXDB_BUCKET")
	token  = os.Getenv("INFLUXDB_TOKEN")
	org    = os.Getenv("INFLUXDB_ORG")
)

var rawInfluxDBClient influxdb2.Client

var (
	productID = "stock.aapl.usa"
	taskID    = "2024-07-05-test"
)

func TestMain(m *testing.M) {
	rawInfluxDBClient = influxdb2.NewClient(url, token)
	m.Run()
	rawInfluxDBClient.Close()
}

func RecreateBucket(client influxdb2.Client, orgName, bucketName string) error {

	org, err := client.OrganizationsAPI().FindOrganizationByName(context.Background(), orgName)
	if err != nil {
		return err
	}

	bucket, err := client.BucketsAPI().FindBucketByName(context.Background(), bucketName)
	if err != nil {
		return nil
	}

	client.BucketsAPI().DeleteBucket(context.Background(), bucket)
	_, err = client.BucketsAPI().CreateBucketWithName(context.Background(), org, bucketName)

	return err
}

func TestOrderEventDispatcher(t *testing.T) {
	t.Run("발송한 order event의 개수와 bucket에 있는 order event의 개수가 같아야 한다.", func(t *testing.T) {
		//arrange
		RecreateBucket(rawInfluxDBClient, org, bucket)
		num := 100
		dispatcher, err := influx.NewOrderEventDispatcher(&influx.Opts{
			Url:        url,
			Token:      token,
			Org:        org,
			BucketName: bucket,
		})

		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		//act
		for i := 0; i < num; i++ {
			dispatcher.Dispatch(taskID, &model.OrderEvent{
				ProductID: productID,
				Command: model.TradeCommand{
					ProportionPercent: 0,
					Action:            model.Buy,
				},
				CreatedAt: time.Now(),
				Task:      model.BackTest,
			})
			// time.Sleep(100 * time.Millisecond)
		}
		dispatcher.Close()
		//assert
		q, err := rawInfluxDBClient.QueryAPI(org).Query(
			context.Background(),
			fmt.Sprintf(`from(bucket:"%s")
				|> range(start:0)
				|> filter(fn: (r)=> r._measurement == "%s")
				|> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
				`, bucket, taskID))
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		count := 0
		for q.Next() {
			assert.Equal(t, productID, q.Record().ValueByKey("productID"))
			assert.Equal(t, int64(0), q.Record().ValueByKey("proportionPercent"))
			assert.Equal(t, "buy", q.Record().ValueByKey("action"))
			assert.Equal(t, "backTest", q.Record().ValueByKey("task"))
			count++
		}
		assert.Equal(t, num, count)
	})
}

func TestAnnotationDispatcher(t *testing.T) {
	// "Testing for the mapper has already been conducted, so specific tests for the mapper will be omitted.
	t.Run("발송한 order event의 개수와 bucket에 있는 order event의 개수가 같아야 한다.", func(t *testing.T) {
		//arrange
		RecreateBucket(rawInfluxDBClient, org, bucket)
		num := 100
		dispatcher, err := influx.NewAnnotationDispatcher(&influx.Opts{
			Url:        url,
			Token:      token,
			Org:        org,
			BucketName: bucket,
		})

		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		type AnnotationSample struct {
			Description string  `name:"description"`
			Price       float64 `name:"price"`
		}
		//act
		for i := 0; i < num; i++ {
			dispatcher.Dispatch(taskID, AnnotationSample{
				Description: "hello world",
				Price:       3.14,
			}, time.Now())
			// time.Sleep(100 * time.Millisecond)
		}
		dispatcher.Close()
		//assert
		q, err := rawInfluxDBClient.QueryAPI(org).Query(
			context.Background(),
			fmt.Sprintf(`from(bucket:"%s")
				|> range(start:0)
				|> filter(fn: (r)=> r._measurement == "%s")
				|> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
				`, bucket, taskID))
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		count := 0
		for q.Next() {
			assert.Equal(t, "hello world", q.Record().ValueByKey("description"))
			assert.Equal(t, float64(3.14), q.Record().ValueByKey("price"))
			count++
		}
		assert.Equal(t, num, count)
	})
}
