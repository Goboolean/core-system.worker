package pipeline_test

import (
	"testing"

	"github.com/Goboolean/core-system.worker/configuration"
	"github.com/Goboolean/core-system.worker/internal/pipeline"
	"github.com/stretchr/testify/assert"
)

func TestBuilder(t *testing.T) {
	t.Run("yml에 가상의 normal pipeline 시나리오가 주어졌을 때, 파이프라인을 적절히 빌드해야 한다.", func(t *testing.T) {
		//arrange
		cfg, err := configuration.ImportAppConfigFromFile("../../test/pipeline_builder_testdata/normal.test.yml")
		if err != nil {
			t.Error(err)
		}

		//act
		p, err := pipeline.Build(*cfg)

		//assert
		assert.NoError(t, err)
		_, ok := p.(*pipeline.Normal)
		assert.True(t, ok)
	})
	t.Run("yml에 가상의 without normal pipeline 시나리오가 주어졌을 때, 파이프라인을 적절히 빌드해야 한다.", func(t *testing.T) {
		//arrange
		cfg, err := configuration.ImportAppConfigFromFile("../../test/pipeline_builder_testdata/without_model.test.yml")
		if err != nil {
			t.Error(err)
		}

		//act
		p, err := pipeline.Build(*cfg)

		//assert
		assert.NoError(t, err)
		_, ok := p.(*pipeline.WithoutModel)
		assert.True(t, ok)
	})
}
