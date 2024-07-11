package fetcher_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/job/fetcher"
	"github.com/Goboolean/core-system.worker/internal/model"
	"github.com/Goboolean/fetch-system.IaC/pkg/influx"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/stretchr/testify/assert"
)

var opts = influx.Opts{
	URL:             os.Getenv("INFLUXDB_URL"),
	Token:           os.Getenv("INFLUXDB_TOKEN"),
	Org:             os.Getenv("INFLUXDB_ORG"),
	TradeBucketName: os.Getenv("INFLUXDB_TRADE_BUCKET"),
}

var rawInfluxClient influxdb2.Client

var testStockID = "stock.aapl.usa"
var testTimeFrame = "1m"

func TestMain(m *testing.M) {
	rawInfluxClient = influxdb2.NewClient(opts.URL, opts.Token)
	m.Run()
	rawInfluxClient.Close()
}

func RecreateBucket(client influxdb2.Client, orgName, bucketName string) error {

	org, err := client.OrganizationsAPI().FindOrganizationByName(context.Background(), orgName)
	if err != nil {
		return err
	}

	bucket, err := client.BucketsAPI().FindBucketByName(context.Background(), bucketName)
	if err != nil {
		return err
	}

	if err := client.BucketsAPI().DeleteBucket(context.Background(), bucket); err != nil {
		return err
	}

	_, err = client.BucketsAPI().CreateBucketWithName(context.Background(), org, bucketName)

	return err
}

func TestPing(t *testing.T) {
	ok, err := rawInfluxClient.Ping(context.Background())
	assert.True(t, ok)
	assert.NoError(t, err)
}

func TestPastStock(t *testing.T) {
	t.Run("저장된 데이터가 없을 때, 0개의 데이터를 가져와야 한다.", func(t *testing.T) {
		//arrange
		if err := RecreateBucket(rawInfluxClient, opts.Org, opts.TradeBucketName); err != nil {
			t.Error(err)
			t.FailNow()
		}
		time.Sleep(100 * time.Millisecond)
		start := time.Now()

		query, err := influx.NewDB(&opts)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second))
		defer cancel()
		err = query.Ping(ctx)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		cursor, err := fetcher.NewStockTradeCursor(query)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		fetchJob, err := fetcher.NewPastStock(cursor, &job.UserParams{
			job.ProductID: testStockID,
			job.StartDate: fmt.Sprint(start.Unix()),
			job.EndDate:   fmt.Sprint(start.Add(time.Minute).Unix()),
			job.TimeFrame: "1m",
		})
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		out := make([]model.Packet, 0)
		go func() {
			for v := range fetchJob.Output() {
				out = append(out, v)
			}
		}()

		err = fetchJob.Execute()

		assert.NoError(t, err)
		assert.Len(t, out, 0)
	})

	t.Run("저장된 데이터가 없을 때, 0개의 데이터를 가져와야 한다.2", func(t *testing.T) {
		//arrange
		if err := RecreateBucket(rawInfluxClient, opts.Org, opts.TradeBucketName); err != nil {
			t.Error(err)
			t.FailNow()
		}
		time.Sleep(100 * time.Millisecond)
		start := time.Now()

		query, err := influx.NewDB(&opts)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second))
		defer cancel()
		err = query.Ping(ctx)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		cursor, err := fetcher.NewStockTradeCursor(query)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		fetchJob, err := fetcher.NewPastStock(cursor, &job.UserParams{
			job.ProductID: testStockID,
			job.StartDate: fmt.Sprint(start.Unix()),
			job.EndDate:   fmt.Sprint(start.Add(time.Minute).Unix()),
			job.TimeFrame: "1m",
		})
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		out := make([]model.Packet, 0)
		go func() {
			for v := range fetchJob.Output() {
				out = append(out, v)
			}
		}()

		err = fetchJob.Execute()

		assert.NoError(t, err)
		assert.Len(t, out, 0)
	})

	// t.Run("데이터가 저장된 만큼, 데이터를 가져와야 한다.", func(t *testing.T) {
	// 	if err := RecreateBucket(rawInfluxClient, opts.Org, opts.TradeBucketName); err != nil {
	// 		t.Error(err)
	// 		t.FailNow()
	// 	}
	// 	time.Sleep(100 * time.Millisecond)
	// 	writer := rawInfluxClient.WriteAPIBlocking(opts.Org, opts.TradeBucketName)
	// 	storeNum := 350
	// 	storeInterval := time.Minute
	// 	start := time.Now().Add(-time.Duration(storeNum) * storeInterval)
	// 	for i := 0; i < storeNum; i++ {
	// 		err := writer.WritePoint(
	// 			context.Background(),
	// 			write.NewPoint(
	// 				fmt.Sprintf("%s.%s", testStockID, testTimeFrame),
	// 				map[string]string{},
	// 				map[string]interface{}{
	// 					"open":   float64(i),
	// 					"close":  float64(2.0),
	// 					"high":   float64(3.0),
	// 					"low":    float64(4.0),
	// 					"volume": float64(4),
	// 				},
	// 				start.Add(time.Duration(i)*storeInterval),
	// 			),
	// 		)

	// 		if err != nil {
	// 			t.Error(err)
	// 			t.FailNow()
	// 		}
	// 	}

	// 	query, err := influx.NewDB(&opts)
	// 	if err != nil {
	// 		t.Error(err)
	// 		t.FailNow()
	// 	}
	// 	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second))
	// 	defer cancel()
	// 	err = query.Ping(ctx)
	// 	if err != nil {
	// 		t.Error(err)
	// 		t.FailNow()
	// 	}

	// 	cursor, err := fetcher.NewStockTradeCursor(query)
	// 	if err != nil {
	// 		t.Error(err)
	// 		t.FailNow()
	// 	}

	// 	fetchJob, err := fetcher.NewPastStock(cursor, &job.UserParams{
	// 		job.ProductID: testStockID,
	// 		job.StartDate: fmt.Sprint(start.Unix()),
	// 		job.EndDate:   fmt.Sprint(start.Add(time.Duration(storeNum) * storeInterval).Unix()),
	// 		job.TimeFrame: "1m",
	// 	})
	// 	if err != nil {
	// 		t.Error(err)
	// 		t.FailNow()
	// 	}

	// 	out := make([]model.Packet, 0)
	// 	go func() {
	// 		for v := range fetchJob.Output() {
	// 			out = append(out, v)
	// 		}
	// 	}()

	// 	err = fetchJob.Execute()

	// 	assert.NoError(t, err)
	// 	assert.Len(t, out, storeNum)

	// })

	t.Run("존재하지 않는 timeFrame일 때 데이터를 가져와선 안 된다.", func(t *testing.T) {
		if err := RecreateBucket(rawInfluxClient, opts.Org, opts.TradeBucketName); err != nil {
			t.Error(err)
			t.FailNow()
		}
		writer := rawInfluxClient.WriteAPIBlocking(opts.Org, opts.TradeBucketName)

		time.Sleep(150 * time.Millisecond)
		storeNum := 350
		storeInterval := time.Minute
		start := time.Now().Add(-time.Duration(storeNum) * storeInterval)
		for i := 0; i < storeNum; i++ {
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
					start.Add(time.Duration(i)*storeInterval),
				),
			)

			if err != nil {
				t.Error(err)
				t.FailNow()
			}
		}

		query, err := influx.NewDB(&opts)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second))
		defer cancel()
		err = query.Ping(ctx)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		cursor, err := fetcher.NewStockTradeCursor(query)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		fetchJob, err := fetcher.NewPastStock(cursor, &job.UserParams{
			job.ProductID: testStockID,
			job.StartDate: fmt.Sprint(start.Unix()),
			job.EndDate:   fmt.Sprint(start.Add(time.Duration(storeNum) * storeInterval).Unix()),
			job.TimeFrame: "1h",
		})
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		out := make([]model.Packet, 0)
		go func() {
			for v := range fetchJob.Output() {
				out = append(out, v)
			}
		}()

		err = fetchJob.Execute()

		assert.NoError(t, err)
		assert.Len(t, out, 0)

	})
}
