package fetcher

import "github.com/Goboolean/core-system.worker/internal/job"

type Fetcher interface {
	job.Common

	Output() job.DataChan
}
