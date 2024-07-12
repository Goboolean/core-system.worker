package model

import (
	"errors"
	"time"
)

// Action represents the action to be performed for an order.
type Action int8

const (
	Buy  Action = -1
	Sell Action = 1
)

// String returns the string representation of the Action.
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

// Task represents the task to be performed by the application.
type Task int8

const (
	BackTest      Task = -1
	RealtimeTrade Task = 1
)

// String returns the string representation of the Task.
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

// ParseTask parses the taskString and returns the corresponding Task.
// If the taskString is invalid, it returns an error.
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

// OrderEvent represents the order event that the worker dispatches to external systems.
type OrderEvent struct {
	ProductID string
	Command   TradeCommand
	CreatedAt time.Time
	Task      Task
}

// TradeCommand represents an order in the system.
type TradeCommand struct {
	ProportionPercent int
	Action            Action
}
