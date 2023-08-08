package model

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	log "github.com/sirupsen/logrus"
)



type ModelSessionImpl struct {
	f *os.File

	ctx context.Context
	cancel context.CancelFunc

	inputChan chan interface{}
	outputChan chan interface{}
}


func newSession(f *os.File) (*ModelSessionImpl, error) {

	ctx, cancel := context.WithCancel(context.Background())

	cmd := exec.CommandContext(ctx, fmt.Sprintf("./%s", f.Name()))

	if err := cmd.Run(); err != nil {
		cancel()
		return nil, err
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		cancel()
		return nil, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		cancel()
		return nil, err
	}

	var instance = ModelSessionImpl{
		f: f,

		ctx: ctx,
		cancel: cancel,		
	}

	// match intputChan with stdin
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case data := <-instance.inputChan:
				fmt.Fprintf(stdin, data.(string))
			}
		}
	}(ctx)

	// match stdout with outputChan
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				buf := make([]byte, 1024)
				n, err := stdout.Read(buf)
				if err != nil {
					if err == io.EOF {
						return
					}
					log.Panic(err)
				}
				fmt.Printf("%s", buf[:n])
			}
		}
	}(ctx)

	return &instance, nil
}


func (m *ModelSessionImpl) GetInputChan() chan<- interface{} {
	return m.inputChan
}

func (m *ModelSessionImpl) GetOutputChan() <-chan interface{} {
	return m.outputChan
}

func (m *ModelSessionImpl) Close() {
	m.cancel()
}
