package job

import "github.com/Goboolean/core-system.worker/internal/model"

type UserParams map[string]string

func (p UserParams) IsKeyNilOrEmpty(key string) bool {
	if val, ok := p[key]; !ok || val == "" {
		return true
	} else {
		return false
	}
}

// User Param Keys
const (
	StartDate = "startDate"
	EndDate   = "endDate"
	ProductID = "productID"
	BatchSize = "batchSize"
	Task      = "task"
	TaskID    = "taskID"
	TimeFrame = "timeFrame"

	NumOfGeneration            = "numOfGeneration"
	MaxRandomDelayMilliseconds = "maxRandomDelayMilliseconds"
)

// DataChan is a channel that is used to send and receive model.Packet objects.
type DataChan chan model.Packet
