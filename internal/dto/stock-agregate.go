package dto

// stockAggreate 우리
type StockAggregate struct {
	OpenTime   int64
	ClosedTime int64
	Open       float64
	Closed     float64
	High       float64
	Low        float64
	Volume     float64
}
