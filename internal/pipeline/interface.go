package pipeline

// Pipeline represents a pipeline that can be run and stopped.
type Pipeline interface {
	Run() error
	Stop()
}
