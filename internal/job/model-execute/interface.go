package modelExecute

import "github.com/Goboolean/core-system.worker/internal/job"

type ModelExecutor interface {
	job.Common

	SetInput(chan any)
	Output() chan any
}
