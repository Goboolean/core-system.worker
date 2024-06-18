package util_test

import (
	"testing"

	"github.com/Goboolean/core-system.worker/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestStopNotifier(t *testing.T) {
	t.Run("Done함수를 호출하지 않았을 때 sn.Done()은 흐름을 block해야 한다.", func(t *testing.T) {

		//arrange
		//act
		sn := util.NewStopNotifier()

		//assert
		assert.NotNil(t, sn)
		assert.NotNil(t, sn.Done())
		select {
		//If a channel returned by sn.Done() is not closed, default will be selected
		case x := <-sn.Done():
			t.Errorf("<-sn.Done() == %v, but it should block", x)
		default:
		}

	})

	t.Run("NotifyStop를 1회 호출했을때 sn.Done은 흐름을 block하지 않아야 한다.", func(t *testing.T) {

		//arrange
		sn := util.NewStopNotifier()
		//act
		sn.NotifyStop()
		//assert
		assert.NotNil(t, sn)
		assert.NotNil(t, sn.Done())

		select {
		//If a channel returned by sn.Done() is closed, the case <-sn.Done() will be selected
		case <-sn.Done():
		default:
			t.Errorf("<-sn.Done() blocked, but it shouldn't block")
		}

	})

	t.Run("NotifyStop를 2회(1회 초과) 호출했을때 sn.Done은 흐름을 block하지 않아야 한다.", func(t *testing.T) {

		//arrange
		sn := util.NewStopNotifier()
		//act
		sn.NotifyStop()
		sn.NotifyStop()
		//assert
		assert.NotNil(t, sn)
		assert.NotNil(t, sn.Done())

		select {
		//If a channel returned by sn.Done() is closed, the case <-sn.Done() will be selected
		case <-sn.Done():
		default:
			t.Errorf("<-sn.Done() blocked, but it shouldn't block")
		}

	})

	t.Run("NotifyStop을 호출했을 때 NotifyStop을 호출하기 전 Done에서 반환한 채널도 흐름을 block하지 않아야 한다.", func(t *testing.T) {

		//arrange
		sn := util.NewStopNotifier()
		done := sn.Done()
		//act
		sn.NotifyStop()
		sn.NotifyStop()
		//assert
		assert.NotNil(t, sn)
		assert.NotNil(t, sn.Done())

		select {
		//If a channel returned by sn.Done() is closed, the case <-sn.Done() will be selected
		case <-done:
		default:
			t.Errorf("<-sn.Done() blocked, but it shouldn't block")
		}

	})
}
