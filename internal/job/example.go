package job

/*
import "context"

type ExampleJob struct {
	Job

	//user param의 type은 float32
	param1 float32

	in  chan any `type:`
	out chan any `type:` //Job은 자신의 Output 채널에 대해 소유권을 가진다.
}

func NewExampleJob(params UserParams) *ExampleJob {
	//여기에 기본값 입력 아웃풋 채널은 job이 소유권을 가져야 한다.
	instance := &ExampleJob{
		out: make(chan any),
	}

	//여기에서 user param 초기화
	if param1, ok := params["param1"]; ok {
		instance.param1 = param1
	}

	return instance
}

func (j *ExampleJob) Execute(ctx context.Context) {
	defer close(j.out)
	go func() {
		select {
		case <-ctx.Done():
			//종료 처리가 왔을 때 처리
		case input, ok := <-j.in:
			if !ok {
				//입력 채널이 닫혔을 때 처리
				return
			}

			data := input.(int)
			//데이터를 이용한 처리
			j.out <- data
		}
	}()

}

func (j *ExampleJob) SetInputChan(input chan any) {
	j.in = input
}

func (j *ExampleJob) OutputChan() chan any {
	return j.out
}

*/
