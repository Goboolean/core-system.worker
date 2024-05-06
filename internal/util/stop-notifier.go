package util

import (
	"sync"
)

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

	if _, ok := <-sn.ch; ok {
		close(sn.ch)
	}
}

func (sn *StopNotifier) Done() chan struct{} {
	return sn.ch
}
