package out

import "context"



type ModelPorts interface {
	NewSession(ctx context.Context, name string) (ModelSession, error)
}

type ModelSession interface {
	// TODO : change interface{} to specific type
	GetInputChan()  chan<- interface{}
	GetOutputChan() <-chan interface{}
	Close() error
}