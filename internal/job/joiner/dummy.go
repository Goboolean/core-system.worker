package joiner

import (
	"github.com/Goboolean/core-system.worker/internal/job"
)

type Dummy struct {
	Joiner

	refIn   job.DataChan
	modelIn job.DataChan
	out     job.DataChan
}
