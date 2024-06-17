package job

// Common is an interface that defines common methods for a job.
type Common interface {
	// Execute executes the job.
	Execute()

	// Close stops job and cleans infra of job and returns an error if any.
	Close() error

	Error() chan error
}
