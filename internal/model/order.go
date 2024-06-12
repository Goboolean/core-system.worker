package model

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
		return "backtest"
	case RealtimeTrade:
		return "realtimetrade"
	default:
		return ""
	}
}

type OrderEvent struct {
	ProductID string
	Command   TradeCommand
	//CreatedAtTimestamp is the timestamp when the order event was created.
	CreatedAtTimestamp int64
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
