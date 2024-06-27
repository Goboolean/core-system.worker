package pipeline_test

import (
	"fmt"
	"testing"

	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/job/analyzer"
	"github.com/Goboolean/core-system.worker/internal/job/fetcher"
	"github.com/Goboolean/core-system.worker/internal/job/transmitter"
	v1 "github.com/Goboolean/core-system.worker/internal/job/transmitter/v1"
	"github.com/Goboolean/core-system.worker/internal/pipeline"
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
		externalChan := make(chan struct{})
		errCh := make(chan error)
		go func() {
			errCh <- p.Run()
		}()

		var stat int
		select {
		case err = <-errCh:
			if err != nil {
				stat = 1
			}
			stat = 0
		case <-externalChan:
			p.Stop()
			stat = 1
		}

		//assert
		assert.NoError(t, err)
		assert.Equal(t, 0, stat)
	})
}
