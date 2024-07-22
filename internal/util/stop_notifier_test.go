package util_test

import (
	"testing"

	"github.com/Goboolean/core-system.worker/internal/util"
	"github.com/stretchr/testify/suite"
)

func TestMain(m *testing.M) {
	m.Run()
}

type StopNotifierTestSuite struct {
	suite.Suite
	sn *util.StopNotifier
}

func (suite *StopNotifierTestSuite) SetupTest() {
	suite.sn = util.NewStopNotifier()
	suite.Require().NotNil(suite.sn)
}

func (suite *StopNotifierTestSuite) TestDoneFunc_ShouldBlock_WhenNotifyStopIsNotCalled() {
	//arrange
	//act
	//assert
	suite.NotNil(suite.sn.Done())
	select {
	//If a channel returned by sn.Done() is not closed, default will be selected
	case x := <-suite.sn.Done():
		suite.Fail("<-sn.Done() == %v, but it should block", x)
	default:
	}

}

func (suite *StopNotifierTestSuite) TestDoneFunc_ShouldNotBlock_WhenNotifyStopIsNotCalledOnce() {
	//arrange
	//act
	suite.sn.NotifyStop()
	//assert

	suite.NotNil(suite.sn.Done())

	select {
	//If a channel returned by sn.Done() is closed, the case <-sn.Done() will be selected
	case <-suite.sn.Done():
	default:
		suite.Fail("<-sn.Done() blocked, but it shouldn't block")
	}

}

func (suite *StopNotifierTestSuite) TestDoneFunc_ShouldNotBlock_WhenNotifyStopIsNotCalledMoreThenOnce() {
	//arrange
	//act
	suite.sn.NotifyStop()
	suite.sn.NotifyStop()

	//assert
	suite.NotNil(suite.sn.Done())

	select {
	//If a channel returned by sn.Done() is closed, the case <-sn.Done() will be selected
	case <-suite.sn.Done():
	default:
		suite.Fail("<-sn.Done() blocked, but it shouldn't block")
	}
}

func (suite *StopNotifierTestSuite) TestChannelRetrievedBeforeCallingNotifyStop_ShouldNotBlock_WhenNotifyStopIsCalled() {
	//arrange
	done := suite.sn.Done()

	//act
	suite.sn.NotifyStop()

	//assert
	suite.NotNil(suite.sn.Done())

	select {
	//If a channel returned by sn.Done() is closed, the case <-sn.Done() will be selected
	case <-done:
	default:
		suite.Fail("<-sn.Done() blocked, but it shouldn't block")
	}

}
