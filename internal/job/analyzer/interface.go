package analyzer

import (
	"github.com/Goboolean/core-system.worker/internal/job"
)

type Analyzer interface {
	job.Common

	SetInput(chan any)
	Output() chan any
}
