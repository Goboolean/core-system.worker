package v1_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/Goboolean/core-system.worker/internal/job"
	"github.com/Goboolean/core-system.worker/internal/job/transmitter"
	v1 "github.com/Goboolean/core-system.worker/internal/job/transmitter/v1"
	"github.com/Goboolean/core-system.worker/internal/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type TestAnnotation struct {
	Number      int
	Description string
}

func TestCommon(t *testing.T) {
	t.Run("annotation과 order가 무작이 순서로 입력됐을 때 적절히 이벤트를 발행해야 한다.", func(t *testing.T) {
		//arrange

		numOrder := 3
		numAnnotation := 3
		inChan := make(job.DataChan, numOrder+numAnnotation)
		fmt.Println("hello")
		i := numOrder
		j := numAnnotation
		for i > 0 || j > 0 {
			switch rand.Intn(2) {
			case 0:
				if i <= 0 {
					continue
				}

				inChan <- model.Packet{
					Sequence: 0,
					Data:     &model.OrderEvent{},
				}
				i--
				fmt.Printf("order event is queued, i:%d\n", i)
			case 1:
				if j <= 0 {
					continue
				}

				inChan <- model.Packet{
					Sequence: 0,
					Data:     &TestAnnotation{},
				}
				j--
				fmt.Printf("annotation is queued, i:%d\n", j)
			}
		}
		close(inChan)
		// arrange mocks
		ctrl := gomock.NewController(t)
		mockOrderEventDispatcher := transmitter.NewMockOrderEventDispatcher(ctrl)
		mockAnnotationDispatcher := transmitter.NewMockAnnotationDispatcher(ctrl)

		mockOrderEventDispatcher.EXPECT().Dispatch(gomock.Any()).Times(numOrder)
		mockAnnotationDispatcher.EXPECT().Dispatch(gomock.Any()).Times(numAnnotation)

		transmit, err := v1.NewCommon(mockAnnotationDispatcher, mockOrderEventDispatcher)
		if err != nil {
			t.Error(err)
		}
		transmit.SetInput(inChan)

		//act
		transmit.Execute()
		for len(inChan) > 0 {
			time.Sleep(250 * time.Millisecond)
		}
		assert.NoError(t, err)
	})
}
