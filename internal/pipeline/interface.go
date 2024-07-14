package pipeline

import "context"

// Represents a interface for pipeline of operations that can be executed.
type Pipeline interface {
	// Executes the pipeline. It takes a context.Context parameter to support cancellation and deadline propagation.
	// It returns an error if any step of the pipeline fails.
	Run(ctx context.Context) error
}
