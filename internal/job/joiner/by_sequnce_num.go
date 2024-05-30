package joiner

import (
	"container/list"
	"sync"

	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/model"
	"github.com/Goboolean/core-system.worker/internal/util"
)

// BySequnceNum는 model.Packet에 있는 sequnce값이 같은 두 데이터를 Pair에 담아 출력합니다.
type BySequnceNum struct {
	Joinner

	refIn   job.DataChan
	modelIn job.DataChan
	out     job.DataChan

	wg   sync.WaitGroup
	stop *util.StopNotifier
}

func NewBysequnce(params *job.UserParams) (*BySequnceNum, error) {

	instance := &BySequnceNum{
		out:  make(job.DataChan),
		wg:   sync.WaitGroup{},
		stop: util.NewStopNotifier(),
	}

	return instance, nil
}

func (b *BySequnceNum) Execute() {
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		defer close(b.out)
		defer b.stop.NotifyStop()

		refanceInputBuf := make([]model.Packet, 0, 100)
		modelInputList := list.New()

		refInChanFail := false
		modelInChanFail := false

		for {
			if refInChanFail && modelInChanFail {
				return
			}

			select {

			case <-b.stop.Done():
				return
			case refrenceDataPacket, ok := <-b.refIn:
				if !ok {
					refInChanFail = true
					continue
				}

				refanceInputBuf = append(refanceInputBuf, refrenceDataPacket)
			case modelDataPacket, ok := <-b.modelIn:
				if !ok {
					modelInChanFail = true
					continue
				}

				modelInputList.PushBack(modelDataPacket)

			}

			for e := modelInputList.Front(); e != nil; e = e.Next() {

				if len(refanceInputBuf) == 0 {
					break
				}

				location := findLargestPacketIndexBySequence(refanceInputBuf, e.Value.(model.Packet).Sequnce)

				if refanceInputBuf[location].Sequnce != e.Value.(model.Packet).Sequnce {
					continue
				}

				b.out <- model.Packet{
					Sequnce: e.Value.(model.Packet).Sequnce,
					Data: &model.Pair{
						RefData:   refanceInputBuf[location].Data,
						ModelData: e.Value.(model.Packet).Data,
					},
				}

				refanceInputBuf = refanceInputBuf[min(len(refanceInputBuf), location+1):]
				modelInputList.Remove(e)
			}

		}
	}()
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

	first := data[0].Sequnce
	last := data[sz-1].Sequnce

	if first > target {
		return -1
	}

	if last <= target {
		return sz - 1
	}

	for i := 0; i < sz-1; i++ {
		if data[i+1].Sequnce > target {
			return i
		}
	}

	return -2
}

func (b *BySequnceNum) Stop() error {
	b.stop.NotifyStop()
	b.wg.Wait()
	return nil
}

func (b *BySequnceNum) SetRefInput(in job.DataChan) {
	b.refIn = in
}

func (b *BySequnceNum) SetModelInput(in job.DataChan) {
	b.modelIn = in
}

func (b *BySequnceNum) Output() job.DataChan {
	return b.out
}
