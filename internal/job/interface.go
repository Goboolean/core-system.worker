package job

// Common is an interface that defines common methods for a job.
type Common interface {

	// Execute runs the given task of the Job.
	// If the Job fails to perform its task, Execute returns an error.
	// If the Job completes successfully, it returns nil.
	// DO NOT CALL Execute() TWICE. IT MUST BE PANIC
	Execute() error
}
