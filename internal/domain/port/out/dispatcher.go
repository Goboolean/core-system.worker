package out

import (
	"github.com/Goboolean/core-system.worker/internal/domain/vo"
)

type WorkDispatcher interface {
	RegisterWorker() <-chan vo.TaskInfo
}
