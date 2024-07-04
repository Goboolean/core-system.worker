package pipeline_test

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
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

		mockOrderEventDispatcher.EXPECT().Dispatch(gomock.Any(), gomock.Any()).Times(num)
		mockOrderEventDispatcher.EXPECT().Close().Times(1)

		mockAnnotationDispatcher.EXPECT().Dispatch(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		mockAnnotationDispatcher.EXPECT().Close().Times(1)

		transmitJob, err := v1.NewCommon(mockAnnotationDispatcher,
			mockOrderEventDispatcher,
			&job.UserParams{})

		if err != nil {
			t.Error(err)
			t.FailNow()
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
			t.FailNow()
		}
		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

		var stat = 0

		externalCh := make(chan struct{})
		done := util.NewStopNotifier()
		go func() {
			select {
			//karfka, message broker
			case <-externalCh:
				cancel()
				stat = 1
			case <-done.Done():
				break
			}
		}()

		err = p.Run(ctx)
		done.NotifyStop()
		if err != nil {
			stat = 1
		}

		//assert
		assert.NoError(t, err)
		assert.Equal(t, 0, stat)
	})
}
