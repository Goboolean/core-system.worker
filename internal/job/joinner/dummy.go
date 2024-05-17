package joinner

type Dummy struct {
	Joinner

	refIn   chan any
	modelIn chan any
	out     chan any
}
