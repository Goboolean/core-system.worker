package checkpatternbasic

import "math"

type Status int

const (
	StatusUndefined Status = iota
	StatusDecreasing
	StatusIncreasing
)

type ChangeStatus int

const (
	ChangeStatusUnchanged ChangeStatus = iota
	ChangeStatusDecreasing
	ChangeStatusIncreasing
)

type TradeEvent int

const (
	TradeEventHold TradeEvent = iota
	TradeEventBuy
	TradeEventSell
)

func (m *Model) calculateChangeStatus(value float64) ChangeStatus {
	if m.prvValue < value && (value-m.currMin) >= m.threshold {
		return ChangeStatusIncreasing
	}

	if m.prvValue > value && (m.currMax-value) >= m.threshold {
		return ChangeStatusDecreasing
	}

	return ChangeStatusUnchanged
}

type Model struct {
	Status    Status
	currMin   float64
	currMax   float64
	prvValue  float64
	threshold float64
}

func NewModel(threshold float64) *Model {
	return &Model{
		Status:    StatusUndefined,
		currMin:   math.MaxFloat64,
		currMax:   -math.MaxFloat64,
		prvValue:  0.0,
		threshold: threshold,
	}
}

func (m *Model) OnEvent(value float64) (float64, TradeEvent) {
	m.currMax = max(m.currMax, value)
	m.currMin = min(m.currMin, value)

	var event TradeEvent = TradeEventHold

	switch m.Status {

	case StatusIncreasing:
		switch m.calculateChangeStatus(value) {
		case ChangeStatusDecreasing:
			m.Status = StatusDecreasing
			event = TradeEventSell
			m.currMin = value
		}

	case StatusDecreasing:
		switch m.calculateChangeStatus(value) {
		case ChangeStatusIncreasing:
			m.Status = StatusIncreasing
			event = TradeEventBuy
			m.currMax = value
		}

	case StatusUndefined:
		switch m.calculateChangeStatus(value) {
		case ChangeStatusIncreasing:
			m.Status = StatusIncreasing
			event = TradeEventBuy
		case ChangeStatusDecreasing:
			m.Status = StatusDecreasing
			event = TradeEventSell
		}
	}

	m.prvValue = value

	var power float64
	if m.currMax == m.currMin {
		power = 0
		return power, event
	}

	switch m.Status {
	case StatusIncreasing:
		power = value - m.currMax
	case StatusDecreasing:
		power = value - m.currMin
	case StatusUndefined:
		power = 0
	}
	return power, event
}
