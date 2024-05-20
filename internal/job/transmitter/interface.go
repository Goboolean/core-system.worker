package transmitter

import (
	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/model"
)

type Transmitter interface {
	job.Common

	SetInput(chan model.Packet)
}
