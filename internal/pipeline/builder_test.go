package pipeline

import (
	"context"
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
	suite.influxC, err = container.InitInfluxContainerWithPortBinding(
		context.Background(),
		"8086",
		configuration.InfluxDBOrderEventBucket,
		configuration.InfluxDBAnnotationBucket)
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
