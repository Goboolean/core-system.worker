package util

import (
	"sync"
)

// closedchan is a reusable closed channel.
var closedchan = make(chan struct{})

func init() {
	close(closedchan)
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

	sn.ch = closedchan
}

func (sn *StopNotifier) Done() chan struct{} {
	return sn.ch
}
