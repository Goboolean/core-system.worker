package fetcher_test

import (
	"strconv"
	"testing"

	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/job/fetcher"
	"github.com/Goboolean/core-system.worker/internal/model"
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

		//act
		out := make([]model.Packet, 0, num)
		outChan := stub.Output()
		stub.Execute()

		for e := range outChan {
			out = append(out, e)
		}
		//assert
		assert.NoError(t, err)
		assert.Len(t, out, num)
	})
}
