package out

import (
	"context"

	"github.com/Goboolean/worker/internal/domain/vo"
)


type ResultDispatcher interface {
	SendResult(context.Context, vo.Result) error
}