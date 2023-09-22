package modelgenerator

import "github.com/Goboolean/core-system.worker/internal/domain/port/out"




type Generator struct {

}

func New(p out.ModelGeneratorPort) (*Generator, error) {
	return &Generator{}, nil
}