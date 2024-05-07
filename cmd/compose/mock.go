package compose

import (
	"context"
	"log"
	"math/rand"
	"time"
)

// It is a function for temporary use for testing kubernetes in advance.
// It will be replaced to another mock class that really mocks other infras.
func MockRun() error {

	min, max := 30, 60
	randSecs := rand.Intn(max-min+1) + min
	randDuration := time.Duration(randSecs) * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), randDuration)
	defer cancel()

	go func() {
		for {
			deadline, _ := ctx.Deadline()
			log.Printf("time left: %s\n", time.Until(deadline))
			time.Sleep(1 * time.Second)
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	}
}
