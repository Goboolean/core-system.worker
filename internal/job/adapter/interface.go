package adapter

import (
	"github.com/Goboolean/core-system.worker/internal/job"
)

type Adapter interface {
	job.Common

	SetInput(chan any)
	Output() chan any
}

// adapt.Common
// adapt.adaptor
