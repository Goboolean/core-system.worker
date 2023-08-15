package model

import (
	"io"
)




type Model struct {
	ch chan struct{}

	stdin io.WriteCloser
	stdout io.ReadCloser
}


