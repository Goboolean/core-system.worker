package job

type JobProvider func(userParams map[string]float32) Job

type JobFactory struct {
	jobRegistry map[string]JobProvider
}

func NewJobFactory() *JobFactory {
	instance := &JobFactory{
		//여기에 Job 입력
	}
	return instance
}

func (f *JobFactory) CreateJob(name string, params UserParams) (Job, error) {

	provider, ok := f.jobRegistry[name]

	if !ok {
		panic("Job is not implied")
	}

	return provider(params), nil
}
