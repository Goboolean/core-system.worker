package model

// Pair is a struct used to pass together
// the output data generated by a machine learning model at a specific point in time (predicted value)
// and the actual trading data at the point(current value).
// It is used as input for analyzers that require both the predicted value and the current value.
type Pair struct {
	ModelData any
	RefData   any
}
