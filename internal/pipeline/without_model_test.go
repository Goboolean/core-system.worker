package pipeline_test

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/job/analyzer"
	"github.com/Goboolean/core-system.worker/internal/job/fetcher"
	"github.com/Goboolean/core-system.worker/internal/job/transmitter"
	v1 "github.com/Goboolean/core-system.worker/internal/job/transmitter/v1"
	"github.com/Goboolean/core-system.worker/internal/pipeline"
	"github.com/Goboolean/core-system.worker/internal/util"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type WithoutModelTestSuite struct {
	suite.Suite
}

func (suite *NormalTestSuite) TestNormal_ShouldFlowDataBetweenJobs_WhenJobInjectedInWithoutModelPipelineWithoutAdapter() {
	//arrange
	num := 100
	fetchJob, err := fetcher.NewStockStub(&job.UserParams{
		"numOfGeneration":            fmt.Sprint(num),
		"maxRandomDelayMilliseconds": fmt.Sprint(5)})
	suite.Require().NoError(err)

	analyzeJob, err := analyzer.NewStub(&job.UserParams{})
	suite.Require().NoError(err)

	ctrl := gomock.NewController(suite.T())
	mockOrderEventDispatcher := transmitter.NewMockOrderEventDispatcher(ctrl)
	mockAnnotationDispatcher := transmitter.NewMockAnnotationDispatcher(ctrl)

	mockOrderEventDispatcher.EXPECT().Dispatch(gomock.Any(), gomock.Any()).Times(num)
	mockOrderEventDispatcher.EXPECT().Close().Times(1)

	mockAnnotationDispatcher.EXPECT().Dispatch(gomock.Any(), gomock.Any(), gomock.Any()).Times(num)
	mockAnnotationDispatcher.EXPECT().Close().Times(1)

	transmitJob, err := v1.NewCommon(mockAnnotationDispatcher,
		mockOrderEventDispatcher,
		&job.UserParams{
			job.TaskID: "2023-3240985",
		})
	suite.Require().NoError(err)

	p, err := pipeline.NewWithoutModelWithoutAdapter(
		fetchJob,
		analyzeJob,
		transmitJob,
	)
	suite.Require().NoError(err)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	var stat = 0

	externalCh := make(chan struct{})
	done := util.NewStopNotifier()
	go func() {
		select {
		//kafka, message broker
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
	suite.NoError(err)
	suite.Equal(0, stat)
}
