package transmitter

import (
	"github.com/Goboolean/core-system.worker/internal/job"
)

type Transmitter interface {
	job.Common

	SetOrderInput(job.DataChan)
	SetAnnotaionInput(job.DataChan)
}
