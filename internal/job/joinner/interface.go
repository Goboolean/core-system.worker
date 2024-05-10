package joinner

import "github.com/Goboolean/core-system.worker/internal/job"

type Joinner interface {
	job.Common

	SetRefInput(chan any)
	SetModelInput(chan any)
	Output() chan any
}
