package out

import (
	"context"

	"github.com/Goboolean/worker/internal/domain/vo"
)



type RealDataFetcher interface {
	GetChannel(context.Context, string) <-chan *vo.StockAggregate
}

type PastDataFetcher interface {
	GetChannel(context.Context, string) <-chan *vo.StockAggregate
}