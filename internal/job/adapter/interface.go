package adapter

import (
	"github.com/Goboolean/core-system.worker/internal/job"
)

type Adapter interface {
	job.Common

	SetInput(job.DataChan)
	Output() job.DataChan
}

// adapt.Common
// adapt.adaptor
