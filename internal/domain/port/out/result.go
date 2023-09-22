package out

import (
	"context"

	"github.com/Goboolean/core-system.worker/internal/domain/vo"
)

type ResultDispatcher interface {
	SendResult(context.Context, *vo.Result) error
}
