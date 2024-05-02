package fetcher_test

import (
	"context"
	"testing"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/core-system.worker/internal/dto"
	"github.com/Goboolean/core-system.worker/internal/infrastructure"
	"github.com/Goboolean/core-system.worker/internal/infrastructure/mock"
	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/job/fetcher"
	"gotest.tools/v3/assert"
)

func TestPastStock(t *testing.T) {
	t.Run("", func(t *testing.T) {

		//Arrange
		mongo, _ := mock.NewMongoClientStock(&resolver.ConfigMap{}, []*infrastructure.StockDocument{
			&infrastructure.StockDocument{
				Symbol:    "AAPL",
				Open:      12,
				Close:     150,
				High:      150,
				Low:       23,
				Average:   0,
				Volume:    12,
				Timestamp: 1711758929,
			},
			&infrastructure.StockDocument{
				Symbol:    "AAPL",
				Open:      12,
				Close:     150,
				High:      150,
				Low:       23,
				Average:   0,
				Volume:    12,
				Timestamp: 1711759229,
			},
			&infrastructure.StockDocument{
				Symbol:    "AAPL",
				Open:      12,
				Close:     150,
				High:      150,
				Low:       23,
				Average:   0,
				Volume:    12,
				Timestamp: 1711759529,
			},
			&infrastructure.StockDocument{
				Symbol:    "AAPL",
				Open:      12,
				Close:     150,
				High:      150,
				Low:       23,
				Average:   0,
				Volume:    12,
				Timestamp: 1711759829,
			},
		})

		ctx, _ := context.WithCancel(context.TODO())

		fetch, err := fetcher.NewPastStock(mongo, &job.UserParams{"timeslice": "1m"})
		if err != nil {
			t.Error(err)
		}

		//Act
		fetch.Execute(ctx)
		out := fetch.Output()
		res := []*dto.StockAggregate{}
		for data := range out {
			val, ok := data.(*dto.StockAggregate)
			if !ok {
				panic("Type miss match")
			}
			res = append(res, val)
		}

		//Assert
		exp := []*dto.StockAggregate{
			&dto.StockAggregate{
				OpenTime:   1711758929,
				ClosedTime: 1711758989,
				Open:       12,
				Closed:     150,
				High:       150,
				Low:        23,
				Volume:     12,
			},
			&dto.StockAggregate{
				OpenTime:   1711759229,
				ClosedTime: 1711759289,
				Open:       12,
				Closed:     150,
				High:       150,
				Low:        23,
				Volume:     12,
			},
			&dto.StockAggregate{
				OpenTime:   1711759529,
				ClosedTime: 1711759589,
				Open:       12,
				Closed:     150,
				High:       150,
				Low:        23,
				Volume:     12,
			},
			&dto.StockAggregate{
				OpenTime:   1711759829,
				ClosedTime: 1711759889,
				Open:       12,
				Closed:     150,
				High:       150,
				Low:        23,
				Volume:     12,
			},
		}
		assert.Equal(t, exp, res)
	})
}
