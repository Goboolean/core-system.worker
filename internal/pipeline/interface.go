package pipeline

// Pipeline represents a pipeline that can be run and stopped.
type Pipeline interface {
	Run()
	Stop()

	// Done returns a channel that is closed when all jobs in the pipeline have completed.
	Done() chan struct{}

	Error() chan error
}
