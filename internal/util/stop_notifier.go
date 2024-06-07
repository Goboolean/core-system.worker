package util

import (
	"sync"
)

// closedChan is a reusable closed channel.
var closedChan = make(chan struct{})

func init() {
	close(closedChan)
}

// StopNotifier is a structure that manages stop.
// It can be used to notify stop to goroutine.
//
// It MUST be instantiated using NewStopNotifier().
type StopNotifier struct {
	mu sync.Mutex

	// The reason for using a channel in this structure is due to the unique properties of channels.
	// When a channel is open, the receiver's code flow is blocked until data is sent to the channel.
	// That is way the channel is not selected in a select-case statement.
	// When a channel is closed, it returns a zero value to the receiver and does not block the receiver's code flow.
	//That is way the channel is selected in a select-case statement.
	ch chan struct{}
}

func NewStopNotifier() *StopNotifier {
	return &StopNotifier{
		ch: make(chan struct{}),
	}
}

// NotifyStop notifies the StopNotifier that a stop signal has been received.
//
// NotifyStop은 2번 이상 호출돼도 예측한 대로 작동합니다.
func (sn *StopNotifier) NotifyStop() {
	sn.mu.Lock()
	defer sn.mu.Unlock()

	// The current channel is changed to closedChan because a closed channel does not block the flow.
	// The reason for assigning closedChan instead of calling the close function is to consider the possibility of multiple calls.
	// Closing an already closed channel causes a panic, so we need to check if the channel is already closed.
	// Since no one is sending data to the channel at this time, this would result in indefinite blocking.
	sn.ch = closedChan
}

// Done returns the channel that will be closed when a stop signal is received.
func (sn *StopNotifier) Done() chan struct{} {
	return sn.ch
}
