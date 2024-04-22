package job

// Common represents a job that can be executed.
type Common interface {
	// Execute executes the job with the given context.
	Execute()
	Close() error
}
