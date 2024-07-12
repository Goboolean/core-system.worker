package model

import "time"

// Packet is a struct used to transfer data between jobs.
// It includes not only data but metadata exchanged between jobs.
type Packet struct {

	// Time indicates the timestamp of the current data, representing the moment the data pertains to.
	// It signifies the time when the target trading data was recorded.
	Time time.Time

	// Data can store any type of data, regardless of its type.
	// If you want to store a struct in Data, it must be of pointer type.
	Data any
}
