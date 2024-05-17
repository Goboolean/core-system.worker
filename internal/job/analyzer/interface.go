package analyzer

import (
	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/model"
)

type Analyzer interface {
	job.Common

	SetInput(chan model.Packet)
	Output() chan model.Packet
}
