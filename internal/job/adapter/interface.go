package adapter

import (
	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/model"
)

type Adapter interface {
	job.Common

	SetInput(chan model.Packet)
	Output() chan model.Packet
}

// adapt.Common
// adapt.adaptor
