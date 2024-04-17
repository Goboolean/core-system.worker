package transmit

import "github.com/Goboolean/core-system.worker/internal/job"

type Transmitter interface {
	job.Common

	SetInput(chan any)
}
