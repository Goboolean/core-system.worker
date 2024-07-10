package integration

import (
	"context"
	"testing"
	"time"

	"github.com/Goboolean/core-system.worker/configuration"
	"github.com/Goboolean/core-system.worker/internal/pipeline"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/stretchr/testify/assert"
)

func TestPipelineWithoutModel(t *testing.T) {
	t.Run("backTesting useCase에서 influxDB에 데이터가 N개 저장됐을 때,"+
		"annotation과 orderEvent에도 N개가 저장돼야 한다.", func(t *testing.T) {
		//arrange
		if err := RecreateBucket(rawInfluxClient, influxDBOrg, tradeBucket); err != nil {
			t.Error(err)
			t.FailNow()
		}
		if err := RecreateBucket(rawInfluxClient, influxDBOrg, orderBucket); err != nil {
			t.Error(err)
			t.FailNow()
		}
		if err := RecreateBucket(rawInfluxClient, influxDBOrg, annotationBucket); err != nil {
			t.Error(err)
			t.FailNow()
		}

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
			if err != nil {
				t.Error(err)
				t.FailNow()
			}
		}

		//act
		config, err := configuration.ImportAppConfigFromFile("./without_model.test.yml")
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		p, err := pipeline.Build(*config)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		ctx := context.Background()
		err = p.Run(ctx)

		//assert
		assert.NoError(t, err)

		var count int
		count, err = CountRecordsInMeasurement(rawInfluxClient, influxDBOrg, tradeBucket, config.TaskID)
		assert.NoError(t, err)
		assert.Equal(t, num, count)

		count, err = CountRecordsInMeasurement(rawInfluxClient, influxDBOrg, orderBucket, config.TaskID)
		assert.NoError(t, err)
		assert.Equal(t, num, count)
	})
}
