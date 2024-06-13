package pipeline

type Pipeline interface {
	Run()
	Stop()

	Done() chan struct{}
}
