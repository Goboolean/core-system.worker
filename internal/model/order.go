package model

type Action int8
type Task int8

const (
	Buy  Action = -1
	Sell Action = 1

	BackTest      Task = -1
	RealtimeTrade Task = 1
)

// TODO: 기술적 요구에 맞게 수정
// 현시점으로 초안
type OrderEvent struct {
	//무엇을 언제 얼마나 사거나 팔 것인가?
	ProductID string
	//목표 비중(퍼센테이지)
	Trade TradeDetails
	// -1이면 즉시 거래?
	// 아니면 현재 unix epoch time?
	// command server는 Timestamp가
	// 그런데 이런 식이면?
	Timestamp int64
}

// Order represents an order in the system.
type TradeDetails struct {
	ProportionPercent int    // Proportion percentage of the order
	Action            Action // Action to be performed for the order
}
