package analyze

import (
	"github.com/Goboolean/core-system.worker/internal/job"
)

type Analyzer interface {
	job.Common

	SetRefInput(chan any)
	SetPredictInput(chan any)
	Output() chan any
}
