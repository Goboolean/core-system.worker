package fetcher_test

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/job/fetcher"
	"github.com/Goboolean/core-system.worker/internal/model"
	"github.com/Goboolean/core-system.worker/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestStub(t *testing.T) {
	// The test takes a long time.
	t.Run("stub에 지정한 개수만큼 output chan에 데이터를 출력해야 한다.", func(t *testing.T) {
		//arrange
		num := 100
		stub, err := fetcher.NewStockStub(&job.UserParams{
			"numOfGeneration":            strconv.FormatInt(int64(num), 10),
			"maxRandomDelayMilliseconds": strconv.FormatInt(100, 10)})
		if err != nil {
			t.Error(err)
			return
		}

		//act

		wg := &sync.WaitGroup{}
		res := make([]model.Packet, 0)

		wg.Add(1)
		go func() {
			defer wg.Done()
			for v := range stub.Output() {
				res = append(res, v)
			}
		}()

		err = stub.Execute()
		if util.IsWaitGroupTimeout(wg, 10*time.Second) {
			t.Errorf("Deadline exceed")
			t.FailNow()
		}

		//assert
		assert.NoError(t, err, 0)
		assert.Len(t, res, num)
	})
}
