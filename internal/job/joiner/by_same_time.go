package joiner

import (
	"container/list"
	"time"

	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/model"
)

// ByTime pairs reference data and model output data that share the same Time
// and passes them as model.Pair objects.
type ByTime struct {
	refIn   job.DataChan
	modelIn job.DataChan
	out     job.DataChan
}

// NewByTime creates new instance of ByTime
func NewByTime(params *job.UserParams) (*ByTime, error) {

	instance := &ByTime{
		out: make(job.DataChan),
	}

	return instance, nil
}

// Execute starts to receive and join data
//
// If the Job fails to perform its task, Execute returns an error.
// If the Job completes successfully, it returns nil.
// DO NOT CALL Execute() TWICE. IT MUST BE PANIC
func (b *ByTime) Execute() error {
	defer close(b.out)

	referenceInputBuf := make([]model.Packet, 0, 100)
	modelInputList := list.New()

	refInChanClosed := false
	modelInChanClosed := false

	for {
		if refInChanClosed && modelInChanClosed {
			return nil
		}

		select {
		case referenceDataPacket, ok := <-b.refIn:
			if !ok {
				refInChanClosed = true
				continue
			}

			referenceInputBuf = append(referenceInputBuf, referenceDataPacket)
		case modelDataPacket, ok := <-b.modelIn:
			if !ok {
				modelInChanClosed = true
				continue
			}

			modelInputList.PushBack(modelDataPacket)

		}

		for e := modelInputList.Front(); e != nil; e = e.Next() {
			if len(referenceInputBuf) == 0 {
				break
			}
			location := findLargestPacketIndexByTime(referenceInputBuf, e.Value.(model.Packet).Time)

			if referenceInputBuf[location].Time != e.Value.(model.Packet).Time {
				continue
			}

			b.out <- model.Packet{
				Time: e.Value.(model.Packet).Time,
				Data: &model.Pair{
					RefData:   referenceInputBuf[location].Data,
					ModelData: e.Value.(model.Packet).Data,
				},
			}

			referenceInputBuf = referenceInputBuf[min(len(referenceInputBuf), location+1):]
			modelInputList.Remove(e)
		}

	}
}

// findLargestPacketIndexBySequence returns the index of the packet with the latest time
// that is before or at the target time.
// -1 means all element has time that is later than target
// WARNING: TO BE USED ONLY WITH ARRAYS SORTED IN ASCENDING ORDER
func findLargestPacketIndexByTime(data []model.Packet, target time.Time) int {
	// data는 순서가 보장돼 있고 대부분 앞 부분에 찾고자 하는 값이 있을 것이라
	// 예상할 수 있으므로 순차탐색
	sz := len(data)
	if sz == 0 {
		return -1
	}

	first := data[0].Time
	last := data[sz-1].Time

	if first.Sub(target) > 0 {
		return -1
	}

	if last.Sub(target) <= 0 {
		return sz - 1
	}

	for i := 0; i < sz-1; i++ {
		if data[i+1].Time.Sub(target) > 0 {
			return i
		}
	}

	return -2
}

func (b *ByTime) SetRefInput(in job.DataChan) {
	b.refIn = in
}

func (b *ByTime) SetModelInput(in job.DataChan) {
	b.modelIn = in
}

func (b *ByTime) Output() job.DataChan {
	return b.out
}
