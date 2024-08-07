package influx_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Goboolean/core-system.worker/internal/infrastructure/influx"
	"github.com/Goboolean/core-system.worker/internal/model"
	influxutil "github.com/Goboolean/core-system.worker/test/util/influx"
	"github.com/stretchr/testify/assert"
)

func TestOrderEventDispatcher(t *testing.T) {
	t.Run("발송한 order event의 개수와 bucket에 있는 order event의 개수가 같아야 한다.", func(t *testing.T) {
		//arrange
		if err := influxutil.RecreateBucket(rawInfluxDBClient, org, orderEventBucket); err != nil {
			t.Error(err)
			t.FailNow()
		}
		num := 100
		dispatcher, err := influx.NewOrderEventDispatcher(&influx.Opts{
			URL:        url,
			Token:      token,
			Org:        org,
			BucketName: orderEventBucket,
		})

		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		//act
		start := time.Now().Add(-time.Duration(num) * time.Second)
		for i := 0; i < num; i++ {
			dispatcher.Dispatch(taskID, &model.OrderEvent{
				ProductID: productID,
				Command: model.TradeCommand{
					ProportionPercent: 0,
					Action:            model.Buy,
				},
				CreatedAt: start.Add(time.Duration(i) * time.Second),
				Task:      model.BackTest,
			})
		}
		dispatcher.Close()
		//assert
		q, err := rawInfluxDBClient.QueryAPI(org).Query(
			context.Background(),
			fmt.Sprintf(`from(bucket:"%s")
				|> range(start:0)
				|> filter(fn: (r)=> r._measurement == "%s")
				|> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
				`, orderEventBucket, taskID))
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
	t.Run("발송한 annotation의 개수와 bucket에 있는 annotation의 개수가 같아야 한다.", func(t *testing.T) {
		//arrange
		if err := influxutil.RecreateBucket(rawInfluxDBClient, org, annotationBucket); err != nil {
			t.Error(err)
			t.FailNow()
		}

		num := 100
		dispatcher, err := influx.NewAnnotationDispatcher(&influx.Opts{
			URL:        url,
			Token:      token,
			Org:        org,
			BucketName: annotationBucket,
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
		start := time.Now().Add(-time.Duration(num) * time.Second)
		for i := 0; i < num; i++ {
			dispatcher.Dispatch(taskID, AnnotationSample{
				Description: "hello world",
				Price:       3.14,
			}, start.Add(time.Duration(i)*time.Second))
		}
		dispatcher.Close()
		//assert
		q, err := rawInfluxDBClient.QueryAPI(org).Query(
			context.Background(),
			fmt.Sprintf(`from(bucket:"%s")
				|> range(start:0)
				|> filter(fn: (r)=> r._measurement == "%s")
				|> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
				`, annotationBucket, taskID))
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
