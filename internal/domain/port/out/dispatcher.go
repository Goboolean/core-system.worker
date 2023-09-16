package out

import (
	"github.com/Goboolean/worker/internal/domain/vo"
)



type WorkDispatcher interface {
	RegisterWorker() <-chan vo.TaskInfo
}