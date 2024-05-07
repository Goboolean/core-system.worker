package fetcher_test

import (
	"testing"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/core-system.worker/internal/dto"
	"github.com/Goboolean/core-system.worker/internal/infrastructure/mongo"
	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/job/fetcher"
	"github.com/stretchr/testify/assert"
)

func TestPastStock(t *testing.T) {
	t.Run("", func(t *testing.T) {

		//Arrange
		mongo, _ := mongo.Mock__NewStockClient(&resolver.ConfigMap{}, []*mongo.StockDocument{
			{
				Symbol:    "AAPL",
				Open:      12,
				Close:     150,
				High:      150,
				Low:       23,
				Average:   0,
				Volume:    12,
				Timestamp: 1711758929,
			}, {
				Symbol:    "AAPL",
				Open:      12,
				Close:     150,
				High:      150,
				Low:       23,
				Average:   0,
				Volume:    12,
				Timestamp: 1711759229,
			}, {
				Symbol:    "AAPL",
				Open:      12,
				Close:     150,
				High:      150,
				Low:       23,
				Average:   0,
				Volume:    12,
				Timestamp: 1711759529,
			}, {
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

		fetch, err := fetcher.NewPastStock(mongo, &job.UserParams{"timeslice": "1m"})
		if err != nil {
			t.Error(err)
		}

		//Act
		fetch.Execute()
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
			{
				OpenTime:   1711758929,
				ClosedTime: 1711758989,
				Open:       12,
				Closed:     150,
				High:       150,
				Low:        23,
				Volume:     12,
			}, {
				OpenTime:   1711759229,
				ClosedTime: 1711759289,
				Open:       12,
				Closed:     150,
				High:       150,
				Low:        23,
				Volume:     12,
			}, {
				OpenTime:   1711759529,
				ClosedTime: 1711759589,
				Open:       12,
				Closed:     150,
				High:       150,
				Low:        23,
				Volume:     12,
			}, {
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
