package model

import (
	"errors"
	"time"
)

type Action int8

const (
	Buy  Action = -1
	Sell Action = 1
)

func (a Action) String() string {
	switch a {
	case Buy:
		return "buy"
	case Sell:
		return "sell"
	default:
		return ""
	}
}

type Task int8

const (
	BackTest      Task = -1
	RealtimeTrade Task = 1
)

func (a Task) String() string {
	switch a {
	case BackTest:
		return "backTest"
	case RealtimeTrade:
		return "realtimeTrade"
	default:
		return ""
	}
}

var ErrInvalidTaskString = errors.New("ParseTask: can't parse taskString")

func ParseTask(taskString string) (Task, error) {
	switch taskString {
	case BackTest.String():
		return BackTest, nil
	case RealtimeTrade.String():
		return RealtimeTrade, nil
	default:
		return 0, ErrInvalidTaskString
	}
}

// TODO: 기술적 요구에 맞게 수정
type OrderEvent struct {
	ProductID string
	Command   TradeCommand
	//CreatedAtTimestamp is the timestamp when the order event was created.
	CreatedAt time.Time
	// Task refers to the operation currently being performed by this application.
	// There are BackTest and RealtimeTrade.
	// External systems can look at this value and
	// decide whether to settle actual trades or simulate trades with past data.
	Task Task
}

// Order represents an order in the system.
type TradeCommand struct {
	ProportionPercent int    // Proportion is the target percentage of the order
	Action            Action // Action to be performed for the order Buy or Sell
}
