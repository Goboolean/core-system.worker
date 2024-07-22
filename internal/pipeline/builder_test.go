package pipeline

import (
	"context"
	"os"
	"testing"

	"github.com/Goboolean/core-system.worker/configuration"
	"github.com/Goboolean/core-system.worker/test/container"
	"github.com/stretchr/testify/suite"
)

type BuildTestSuite struct {
	suite.Suite
	influxC *container.InfluxContainer
}

func (suite *BuildTestSuite) SetupSuite() {
	var err error
	suite.influxC, err = container.InitInfluxContainer(context.Background(),
		os.Getenv("INFLUXDB_TRADE_BUCKET"),
		os.Getenv("INFLUXDB_ORDER_EVENT_BUCKET"),
		os.Getenv("INFLUXDB_ANNOTATION_BUCKET"))

	// 고언어 테스트가 병렬적으로 실행될 때는 멀티프로세스 환경에서 실행되므로
	// 한 패키지에서 설정한 환경변수는 다른 패키지 테스트에서 영향을 미치지 않는다.
	os.Setenv("INFLUXDB_URL", suite.influxC.URL)
	suite.Require().NoError(err)
}

func (suite *BuildTestSuite) TearDownSuite() {
	suite.Require().NoError(
		suite.influxC.Terminate(context.Background()))
}

func (suite *BuildTestSuite) TestBuild_ShouldBuildNormalPipeline_WhenGivenVirtualNormalPipelineScenarioInYMLConfiguration() {
	// arrange
	cfg, err := configuration.ImportAppConfigFromFile("../../test/pipeline_builder_testdata/normal.test.yml")
	suite.Require().NoError(err)

	// act
	p, err := Build(*cfg)

	// assert
	suite.NoError(err)

	_, ok := p.(*Normal)
	suite.True(ok)
}

func (suite *BuildTestSuite) TestBuild_ShouldBuildWithoutNormalPipeline_WhenGivenVirtualNormalPipelineScenarioInYMLConfiguration() {
	//arrange
	cfg, err := configuration.ImportAppConfigFromFile("../../test/pipeline_builder_testdata/without_model.test.yml")
	suite.Require().NoError(err)

	//act
	p, err := Build(*cfg)

	//assert
	suite.NoError(err)
	_, ok := p.(*WithoutModel)
	suite.True(ok)

}

func TestBuilder(t *testing.T) {
	suite.Run(t, new(BuildTestSuite))
}
