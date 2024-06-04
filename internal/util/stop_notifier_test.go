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
		case x := <-sn.Done():
			t.Errorf("<-sn.Done() == %v, but it should block", x)
		default:
		}

	})

	t.Run("NotifyStop를 1회 호출했을때 sn.Done은 흐름을 block해야 한다.", func(t *testing.T) {

		//arrange
		sn := util.NewStopNotifier()
		//act
		sn.NotifyStop()
		//assert
		assert.NotNil(t, sn)
		assert.NotNil(t, sn.Done())

		select {
		case <-sn.Done():
		default:
			t.Errorf("<-sn.Done() blocked, but it shouldn't block")
		}

	})

	t.Run("NotifyStop를 2회(1회 초과) 호출했을때 sn.Done은 흐름을 block해야 한다.", func(t *testing.T) {

		//arrange
		sn := util.NewStopNotifier()
		//act
		sn.NotifyStop()
		sn.NotifyStop()
		//assert
		assert.NotNil(t, sn)
		assert.NotNil(t, sn.Done())

		select {
		case <-sn.Done():
		default:
			t.Errorf("<-sn.Done() blocked, but it shouldn't block")
		}

	})

}
