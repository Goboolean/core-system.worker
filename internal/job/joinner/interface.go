package joinner

import (
	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/model"
)

type Joinner interface {
	job.Common

	SetRefInput(chan model.Packet)
	SetModelInput(chan model.Packet)
	Output() chan model.Packet
}
