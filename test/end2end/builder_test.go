package end2end

import (
	"testing"

	"github.com/Goboolean/core-system.worker/configuration"
	"github.com/Goboolean/core-system.worker/internal/pipeline"
	"github.com/stretchr/testify/suite"
)

type BuildTestSuite struct {
	suite.Suite
}

func (suite *BuildTestSuite) TestBuild_ShouldBuildNormalPipeline_WhenGivenVirtualNormalPipelineScenarioInYMLConfiguration() {
	// arrange
	cfg, err := configuration.ImportAppConfigFromFile("../../test/pipeline_builder_testdata/normal.test.yml")
	suite.Require().NoError(err)

	// act
	p, err := pipeline.Build(*cfg)

	// assert
	suite.NoError(err)

	_, ok := p.(*pipeline.Normal)
	suite.True(ok)
}

func (suite *BuildTestSuite) TestBuild_ShouldBuildWithoutNormalPipeline_WhenGivenVirtualNormalPipelineScenarioInYMLConfiguration() {
	//arrange
	cfg, err := configuration.ImportAppConfigFromFile("../../test/pipeline_builder_testdata/without_model.test.yml")
	suite.Require().NoError(err)

	//act
	p, err := pipeline.Build(*cfg)

	//assert
	suite.NoError(err)
	_, ok := p.(*pipeline.WithoutModel)
	suite.True(ok)

}

func TestBuilder(t *testing.T) {
	suite.Run(t, new(BuildTestSuite))
}
