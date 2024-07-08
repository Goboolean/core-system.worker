package model

// ExampleAnnotation is the most basic model of Annotation,
// storing human-readable interpretation results at a specific point in time.
type ExampleAnnotation struct {
	Description string `name:"description"`
}
