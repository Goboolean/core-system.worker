package pipeline

import (
	"context"
	"errors"

	"github.com/Goboolean/core-system.worker/internal/job/adapter"
	"github.com/Goboolean/core-system.worker/internal/job/analyzer"
	"github.com/Goboolean/core-system.worker/internal/job/executer"
	"github.com/Goboolean/core-system.worker/internal/job/fetcher"
	"github.com/Goboolean/core-system.worker/internal/job/joiner"
	"github.com/Goboolean/core-system.worker/internal/job/transmitter"
	"github.com/Goboolean/core-system.worker/internal/model"
	"github.com/Goboolean/core-system.worker/internal/util"
	"github.com/Goboolean/core-system.worker/internal/util/chanutil"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

var ErrTypeNotMatch = errors.New("pipeline: cannot build a pipeline because the types are not compatible between the jobs")

// 아키텍처 설계 상 이 구조는 변경되면 안 된다.
type Normal struct {
	//jobs
	fetcher       fetcher.Fetcher
	joiner        joiner.Joiner
	modelExecuter executer.ModelExecutor
	adapter       adapter.Adapter
	resAnalyzer   analyzer.Analyzer
	transmitter   transmitter.Transmitter

	//utils
	mux  *chanutil.ChannelMux[model.Packet]
	done *util.StopNotifier
}

func NewNormalWithAdapter(
	fetcher fetcher.Fetcher,
	joiner joiner.Joiner,
	modelExecuter executer.ModelExecutor,
	adapter adapter.Adapter,
	resAnalyzer analyzer.Analyzer,
	transmitter transmitter.Transmitter) (*Normal, error) {

	instance := Normal{
		fetcher:       fetcher,
		joiner:        joiner,
		modelExecuter: modelExecuter,
		adapter:       adapter,
		resAnalyzer:   resAnalyzer,
		transmitter:   transmitter,

		mux:  chanutil.NewChannelMux[model.Packet](),
		done: util.NewStopNotifier(),
	}

	instance.mux.SetInput(instance.fetcher.Output())
	instance.modelExecuter.SetInput(instance.mux.Output())
	instance.adapter.SetInput(instance.modelExecuter.Output())
	instance.joiner.SetModelInput(instance.adapter.Output())
	instance.joiner.SetRefInput(instance.mux.Output())
	instance.resAnalyzer.SetInput(instance.joiner.Output())
	instance.transmitter.SetInput(instance.resAnalyzer.Output())

	return &instance, nil
}

func NewNormalWithoutAdapter(
	fetch fetcher.Fetcher,
	join joiner.Joiner,
	modelExec executer.ModelExecutor,
	resAnalyze analyzer.Analyzer,
	transmit transmitter.Transmitter) (*Normal, error) {

	instance := Normal{
		fetcher:       fetch,
		joiner:        join,
		modelExecuter: modelExec,
		resAnalyzer:   resAnalyze,
		transmitter:   transmit,

		mux:  chanutil.NewChannelMux[model.Packet](),
		done: util.NewStopNotifier(),
	}

	instance.mux.SetInput(instance.fetcher.Output())
	instance.modelExecuter.SetInput(instance.mux.Output())
	instance.joiner.SetModelInput(instance.modelExecuter.Output())
	instance.joiner.SetRefInput(instance.mux.Output())
	instance.resAnalyzer.SetInput(instance.joiner.Output())
	instance.transmitter.SetInput(instance.resAnalyzer.Output())

	return &instance, nil

}

func (n *Normal) Run(ctx context.Context) error {
	g := errgroup.Group{}
	stop := util.StopNotifier{}

	go func() {
		select {
		case <-stop.Done():
			n.fetcher.NotifyStop()
			break
		case <-ctx.Done():
			n.fetcher.NotifyStop()
			break
		case <-n.done.Done():
			break
		}
	}()

	n.mux.Execute()

	g.Go(func() error {
		err := n.fetcher.Execute()
		log.Debug("fetch job is completed")
		return err
	})

	g.Go(func() error {
		err := n.joiner.Execute()
		if err != nil {
			stop.NotifyStop()
		}
		log.Debug("join job is completed")
		return err
	})

	g.Go(func() error {
		err := n.modelExecuter.Execute()
		if err != nil {
			stop.NotifyStop()
		}
		log.Debug("execute job is completed")
		return err
	})

	g.Go(func() error {
		if n.adapter == nil {
			return nil
		}

		err := n.adapter.Execute()
		if err != nil {
			stop.NotifyStop()
		}
		log.Debug("adapt job is completed")
		return err
	})

	g.Go(func() error {
		err := n.resAnalyzer.Execute()
		if err != nil {
			stop.NotifyStop()
		}
		log.Debug("analyze job is completed")
		return err
	})

	g.Go(func() error {
		err := n.transmitter.Execute()
		if err != nil {
			stop.NotifyStop()
		}
		log.Debug("transmit job is completed")
		return err
	})

	var err error
	go func() {
		err = g.Wait()
		log.Info("Pipeline job is completed")
		n.done.NotifyStop()
	}()

	<-n.done.Done()
	return err
}
