package joinner

import (
	"github.com/Goboolean/core-system.worker/internal/job"
)

type Dummy struct {
	Joinner

	refIn   job.DataChan
	modelIn job.DataChan
	out     job.DataChan
}
