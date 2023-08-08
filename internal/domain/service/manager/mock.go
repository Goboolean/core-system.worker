package manager

import (
	"math/rand"
	"time"
)

type Mock struct {

}


func NewMock() *Mock {
	return &Mock{}
}



func (m *Mock) Run() error {
	rand.Seed(time.Now().UnixNano())

	min, max := 30, 60
	randSecs := rand.Intn(max-min+1) + min
	randDuration := time.Duration(randSecs) * time.Second

	time.Sleep(randDuration)
	return nil
}