package util

import (
	"sync"
)

// StopNotifier is a structure that manages stop.
// It can be used to notify stop to goroutine.
//
// It MUST be instantiated using NewStopNotifier().
type StopNotifier struct {
	mu sync.Mutex
	//once는 NotifyStop이 2번 이상 호출돼도 ch를 오직 한 번만 close하도록 합니다.
	once sync.Once

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

	closeChan := func() {
		close(sn.ch)
	}
	sn.once.Do(closeChan)
}

// Done returns the channel that will be closed when a stop signal is received.
func (sn *StopNotifier) Done() chan struct{} {
	return sn.ch
}
