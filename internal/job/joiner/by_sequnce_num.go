package joiner

import (
	"container/list"

	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/model"
)

// BySequenceNum는 model.Packet에 있는 sequnce값이 같은 두 데이터를 Pair에 담아 출력합니다.
type BySequenceNum struct {
	refIn   job.DataChan
	modelIn job.DataChan
	out     job.DataChan
}

func NewBySequence(params *job.UserParams) (*BySequenceNum, error) {

	instance := &BySequenceNum{
		out: make(job.DataChan),
	}

	return instance, nil
}

func (b *BySequenceNum) Execute() error {
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
			location := findLargestPacketIndexBySequence(referenceInputBuf, e.Value.(model.Packet).Sequence)

			if referenceInputBuf[location].Sequence != e.Value.(model.Packet).Sequence {
				continue
			}

			b.out <- model.Packet{
				Sequence: e.Value.(model.Packet).Sequence,
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

// findLargestPacketIndexBySequence returns the index of the packet with the largest sequence number
// that is less than or equal to the target sequence number.
// -1 means all element has sequnce that is grater than target
// WARINING: TO BE USED ONLY WITH ARRAYS SORTED IN ASCENDING ORDER
func findLargestPacketIndexBySequence(data []model.Packet, target int64) int {
	// data는 순서가 보장돼 있고 대부분 앞 부분에 찾고자 하는 값이 있을 것이라
	// 예상할 수 있으므로 순차탐색
	sz := len(data)
	if sz == 0 {
		return -1
	}

	first := data[0].Sequence
	last := data[sz-1].Sequence

	if first > target {
		return -1
	}

	if last <= target {
		return sz - 1
	}

	for i := 0; i < sz-1; i++ {
		if data[i+1].Sequence > target {
			return i
		}
	}

	return -2
}

func (b *BySequenceNum) SetRefInput(in job.DataChan) {
	b.refIn = in
}

func (b *BySequenceNum) SetModelInput(in job.DataChan) {
	b.modelIn = in
}

func (b *BySequenceNum) Output() job.DataChan {
	return b.out
}
