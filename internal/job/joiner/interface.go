package joiner

import (
	"github.com/Goboolean/core-system.worker/internal/job"
)

type Joinner interface {
	job.Common

	SetRefInput(job.DataChan)
	SetModelInput(job.DataChan)
	Output() job.DataChan
}
