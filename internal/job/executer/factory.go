package executer

import (
	"fmt"

	"github.com/Goboolean/core-system.worker/internal/job"
)

type jobProvider func(p *job.UserParams) (ModelExecutor, error)

// wire의 한계로 여기서 수동DI 합니다.
var providerRepo = map[Spec]jobProvider{}

func CreateExecuter(spec Spec, p *job.UserParams) (ModelExecutor, error) {

	var provider, ok = providerRepo[spec]
	if !ok {
		return nil, job.NotFoundJob
	}

	f, err := provider(p)
	if err != nil {
		return nil, fmt.Errorf("create model execute job: %w", err)
	}

	return f, nil
}
