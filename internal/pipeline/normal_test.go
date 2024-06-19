package pipeline_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/job/analyzer"
	"github.com/Goboolean/core-system.worker/internal/job/executer"
	"github.com/Goboolean/core-system.worker/internal/job/fetcher"
	"github.com/Goboolean/core-system.worker/internal/job/joiner"
	"github.com/Goboolean/core-system.worker/internal/job/transmitter"
	v1 "github.com/Goboolean/core-system.worker/internal/job/transmitter/v1"
	"github.com/Goboolean/core-system.worker/internal/pipeline"
	"github.com/Goboolean/core-system.worker/internal/util"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNormal(t *testing.T) {
	t.Run("어뎁터가 필요하지 않은 normal pipeline에 job을 주입했을 때 job사이에서 데이터가 흘러야 한다.", func(t *testing.T) {
		//arrange
		num := 100
		fetchJob, err := fetcher.NewStockStub(&job.UserParams{
			"numOfGeneration":            fmt.Sprint(num),
			"maxRandomDelayMilliseconds": fmt.Sprint(5)})
		if err != nil {
			t.Error(err)
			return
		}
		executeJob, err := executer.NewStub(&job.UserParams{})
		if err != nil {
			t.Error(err)
			return
		}
		analyzeJob, err := analyzer.NewStub(&job.UserParams{})
		if err != nil {
			t.Error(err)
			return
		}
		joinJob, err := joiner.NewBySequence(&job.UserParams{})
		if err != nil {
			t.Error(err)
			return
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
			return
		}

		p, err := pipeline.NewNormalWithoutAdapter(
			fetchJob,
			joinJob,
			executeJob,
			analyzeJob,
			transmitJob,
		)

		if err != nil {
			t.Error(err)
			return
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
		}
		wg.Wait()

		//assert
		assert.Len(t, errInPipeline, 0)
		assert.Equal(t, 0, stat)
	})
}
