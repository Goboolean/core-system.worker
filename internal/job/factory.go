package job

type UserParams map[string]string

func (p UserParams) IsKeyNullOrEmpty(key string) bool {
	if val, ok := p[key]; !ok || val == "" {
		return true
	} else {
		return false
	}
}

type JobProvider func(userParams *UserParams) Job
type JobFactory struct {
	jobRegistry map[string]JobProvider
}

func NewJobFactory() *JobFactory {
	instance := &JobFactory{
		//여기에 Job 입력
	}
	return instance
}

func (f *JobFactory) CreateJob(name string, params *UserParams) (Job, error) {

	provider, ok := f.jobRegistry[name]

	if !ok {
		panic("Job is not implied")
	}

	return provider(params), nil
}
