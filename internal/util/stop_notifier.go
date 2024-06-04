package util

import (
	"sync"
)

// closedChan is a reusable closed channel.
var closedChan = make(chan struct{})

func init() {
	close(closedChan)
}

type StopNotifier struct {
	mu sync.Mutex
	ch chan struct{}
}

func NewStopNotifier() *StopNotifier {
	return &StopNotifier{
		ch: make(chan struct{}),
	}
}

func (sn *StopNotifier) NotifyStop() {
	sn.mu.Lock()
	defer sn.mu.Unlock()

	sn.ch = closedChan
}

func (sn *StopNotifier) Done() chan struct{} {
	return sn.ch
}
