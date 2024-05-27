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

				location := findLargestPacketIndexBySequence(refanceInputBuf, e.Value.(model.Packet).Sequnce)

				if refanceInputBuf[location].Sequnce != e.Value.(model.Packet).Sequnce {
					break
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
func findLargestPacketIndexBySequence(data []model.Packet, target int64) int {
	size := len(data)

	i := 0
	for i < size && data[i].Sequnce < target {
		i++
	}

	return i
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
