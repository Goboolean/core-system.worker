package adapt

import (
	"github.com/Goboolean/core-system.worker/internal/job"
)

type Adapter interface {
	job.Common

	SetRefInput(chan any)
	SetPredictInput(chan any)
	Output() chan any
}

// adapt.Common
// adapt.adaptor
