package pipeline

import "context"

// Pipeline represents a pipeline that can be run and stopped.
type Pipeline interface {
	Run(ctx context.Context) error
}
