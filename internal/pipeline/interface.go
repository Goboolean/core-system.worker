package pipeline

type Pipeline interface {
	Run()
	Stop() error
}
