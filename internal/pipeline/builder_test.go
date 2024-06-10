package pipeline_test

import (
	"reflect"
	"testing"

	"github.com/Goboolean/core-system.worker/configuration"
	"github.com/Goboolean/core-system.worker/internal/pipeline"
	"github.com/stretchr/testify/assert"
)

func TestBuilder(t *testing.T) {
	t.Run("yml에 가상의 normal pipeline 시나리오가 주어졌을 때, 파이프라인을 적절히 빌드해야 한다.", func(t *testing.T) {
		//arrange
		cfg, err := configuration.ImportAppConfigFromFile("./testdata/normal.test.yml")
		if err != nil {
			t.Error(err)
		}

		p, err := pipeline.Build(*cfg)
		assert.NoError(t, err)
		assert.Equal(t, "Normal", reflect.TypeOf(p).Name())
	})
	t.Run("yml에 가상의 without normal pipeline 시나리오가 주어졌을 때, 파이프라인을 적절히 빌드해야 한다.", func(t *testing.T) {
		//arrange
		cfg, err := configuration.ImportAppConfigFromFile("without_model.test.yml")
		if err != nil {
			t.Error(err)
		}

		p, err := pipeline.Build(*cfg)
		if err != nil {
			t.Error(err)
		}

		//act
		p.Run()
		//assert
		//??
	})
}
