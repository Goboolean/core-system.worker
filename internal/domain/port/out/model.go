package out

import (
	"context"

	"github.com/Goboolean/core-system.worker/internal/domain/vo"
)

type ModelGeneratorPort interface {
	
}


type ModelPort interface {
	NewSession(ctx context.Context, name string) (ModelSession, error)
}

type ModelSession interface {
	// TODO : change interface{} to specific type
	GetInputChan()  chan<- *vo.StockAggregate
	GetOutputChan() <-chan *vo.Result
	Close() error
}