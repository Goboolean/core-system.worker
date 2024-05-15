package executer

import "github.com/Goboolean/core-system.worker/internal/job"

type ModelExecutor interface {
	job.Common

	SetInput(job.DataChan)
	Output() job.DataChan
}
