package pipeline

import (
	"context"

	"github.com/Goboolean/core-system.worker/internal/job/adapter"
	"github.com/Goboolean/core-system.worker/internal/job/analyzer"
	"github.com/Goboolean/core-system.worker/internal/job/fetcher"
	"github.com/Goboolean/core-system.worker/internal/job/transmitter"
	"github.com/Goboolean/core-system.worker/internal/util"
	"golang.org/x/sync/errgroup"
)

type WithoutModel struct {
	fetcher     fetcher.Fetcher
	adapter     adapter.Adapter
	analyzer    analyzer.Analyzer
	transmitter transmitter.Transmitter

	done *util.StopNotifier
}

func NewWithoutModelWithAdapter(
	fetcher fetcher.Fetcher,
	adapter adapter.Adapter,
	analyzer analyzer.Analyzer,
	transmitter transmitter.Transmitter) (*WithoutModel, error) {

	instance := WithoutModel{
		fetcher:     fetcher,
		adapter:     adapter,
		analyzer:    analyzer,
		transmitter: transmitter,

		done: util.NewStopNotifier(),
	}

	instance.adapter.SetInput(instance.fetcher.Output())
	instance.analyzer.SetInput(instance.adapter.Output())
	instance.transmitter.SetInput(instance.analyzer.Output())

	return &instance, nil
}

func NewWithoutModelWithoutAdapter(
	fetch fetcher.Fetcher,
	analyze analyzer.Analyzer,
	transmit transmitter.Transmitter) (*WithoutModel, error) {

	instance := WithoutModel{
		fetcher:     fetch,
		analyzer:    analyze,
		transmitter: transmit,
		done:        util.NewStopNotifier(),
	}

	instance.analyzer.SetInput(instance.fetcher.Output())
	instance.transmitter.SetInput(instance.analyzer.Output())

	return &instance, nil
}

func (wom *WithoutModel) Run(ctx context.Context) error {
	g := errgroup.Group{}
	stop := util.StopNotifier{}
	go func() {
		select {
		case <-stop.Done():
			wom.fetcher.NotifyStop()
			break
		case <-ctx.Done():
			wom.fetcher.NotifyStop()
			break
		case <-wom.done.Done():
			break
		}
	}()

	g.Go(func() error {
		return wom.fetcher.Execute()
	})

	g.Go(func() error {
		if wom.adapter == nil {
			return nil
		}

		err := wom.adapter.Execute()
		if err != nil {
			stop.NotifyStop()
		}
		return err
	})

	g.Go(func() error {
		err := wom.analyzer.Execute()
		if err != nil {
			stop.NotifyStop()
		}
		return err
	})

	g.Go(func() error {
		err := wom.transmitter.Execute()
		if err != nil {
			stop.NotifyStop()
		}
		return err
	})

	var err error
	go func() {
		err = g.Wait()
		wom.done.NotifyStop()
	}()

	<-wom.done.Done()
	return err
}
