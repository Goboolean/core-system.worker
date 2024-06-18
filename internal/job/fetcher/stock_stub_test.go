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
	//테스트에 시간이 오래 걸립니다.
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
		stub.Execute()
		res := make([]model.Packet, 0)
		errInJob := make([]error, 0)

		wg.Add(1)
		go func() {
			defer wg.Done()
			for v := range stub.Output() {
				res = append(res, v)
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			for v := range stub.Error() {
				errInJob = append(errInJob, v)
			}
		}()

		if util.IsWaitGroupTimeout(wg, 10*time.Second) {
			t.Errorf("Deadline exceed")
		}

		//asse
		assert.Len(t, res, num)
		assert.Len(t, errInJob, 0)
	})
}
