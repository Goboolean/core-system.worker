package model

// Packet is a struct used to transfer data between jobs.
// It includes not only data but metadata exchanged between jobs.
type Packet struct {
	// Sequence represents the order in which the data was generated.
	// The sequence remains unchanged when moving between jobs.
	Sequence int64

	// Data can store any type of data, regardless of its type.
	// If you want to store a struct in Data, it must be of pointer type.
	Data any
}
