package pipeline_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/job/analyzer"
	"github.com/Goboolean/core-system.worker/internal/job/fetcher"
	"github.com/Goboolean/core-system.worker/internal/job/transmitter"
	v1 "github.com/Goboolean/core-system.worker/internal/job/transmitter/v1"
	"github.com/Goboolean/core-system.worker/internal/pipeline"
	"github.com/Goboolean/core-system.worker/internal/util"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestWithoutMode(t *testing.T) {
	t.Run("어뎁터가 필요하지 않은 without model pipeline에 job을 주입했을 때 job사이에서 데이터가 흘러야 한다.", func(t *testing.T) {
		//arrange
		num := 100
		fetchJob, err := fetcher.NewStockStub(&job.UserParams{
			"numOfGeneration":            fmt.Sprint(num),
			"maxRandomDelayMilliseconds": fmt.Sprint(5)})
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		analyzeJob, err := analyzer.NewStub(&job.UserParams{})
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		ctrl := gomock.NewController(t)
		mockOrderEventDispatcher := transmitter.NewMockOrderEventDispatcher(ctrl)
		mockAnnotationDispatcher := transmitter.NewMockAnnotationDispatcher(ctrl)

		mockOrderEventDispatcher.EXPECT().Dispatch(gomock.Any()).Times(num)
		mockOrderEventDispatcher.EXPECT().Close().Times(1)

		mockAnnotationDispatcher.EXPECT().Dispatch(gomock.Any()).AnyTimes()
		mockAnnotationDispatcher.EXPECT().Close().Times(1)

		transmitJob, err := v1.NewCommon(mockAnnotationDispatcher,
			mockOrderEventDispatcher,
			&job.UserParams{})

		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		p, err := pipeline.NewWithoutModelWithoutAdapter(
			fetchJob,
			analyzeJob,
			transmitJob,
		)

		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		//act
		errInPipeline := make([]error, 0)
		p.Run()

		externalChan := make(chan struct{})

		stop := util.NewStopNotifier()
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			for e := range p.Error() {
				errInPipeline = append(errInPipeline, e)
				stop.NotifyStop()
			}
		}()

		var stat int

		select {
		case <-p.Done():
			stat = 0
		case <-externalChan:
			p.Stop()
			stat = 1
		case <-stop.Done():
			p.Stop()
			stat = 1
		case <-time.After(2 * time.Second):
			t.FailNow()
		}
		wg.Wait()

		//assert
		assert.Len(t, errInPipeline, 0)
		assert.Equal(t, 0, stat)
	})
}