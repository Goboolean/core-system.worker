package vo

type TaskType int64

const (
	Real TaskType = iota+1
	Past
)

type TaskInfo struct {
	Type TaskType
	StockId string
	ModelId string
}