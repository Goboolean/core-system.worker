package fetcher

import (
	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/model"
)

type Fetcher interface {
	job.Common

	Output() chan model.Packet
}
