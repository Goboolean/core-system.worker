package executer

import (
	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/model"
)

type ModelExecutor interface {
	job.Common

	SetInput(chan model.Packet)
	Output() chan model.Packet
}
