package checkpatternbasic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXxx(t *testing.T) {
	type outputFormat struct {
		rate  float64
		event TradeEvent
	}

	// Arrange
	var tests = []struct {
		model  *Model
		input  []float64
		output []outputFormat
	}{
		{
			NewModel(3),
			[]float64{0, 1, 2, 4, 5},
			[]outputFormat{
				{0, TradeEventHold},
				{0, TradeEventHold},
				{0, TradeEventHold},
				{0, TradeEventBuy},
				{0, TradeEventHold},
			},
		},
		{
			NewModel(2.5),
			[]float64{5, 4, 3, 2, 1},
			[]outputFormat{
				{0, TradeEventHold},
				{0, TradeEventHold},
				{0, TradeEventHold},
				{0, TradeEventSell},
				{0, TradeEventHold},
			},
		},
		{
			NewModel(3),
			[]float64{1, 3, 5, 4, 2},
			[]outputFormat{
				{0, TradeEventHold},
				{0, TradeEventHold},
				{0, TradeEventBuy},
				{-1, TradeEventHold},
				{0, TradeEventSell},
			},
		},
		{
			NewModel(5),
			[]float64{1, 3, 7, 5, 2, 3, 4, 1},
			[]outputFormat{
				{0, TradeEventHold},
				{0, TradeEventHold},
				{0, TradeEventBuy},
				{-2, TradeEventHold},
				{0, TradeEventSell},
				{1, TradeEventHold},
				{2, TradeEventHold},
				{0, TradeEventHold},
			},
		},
	}

	for i, v := range tests {
		for j, value := range v.input {
			// Act
			rate, event := v.model.OnEvent(value)

			// Assert
			assert.Equal(t, v.output[j].rate, rate, "Test case %d, input %d", i, j)
			assert.Equal(t, v.output[j].event, event, "Test case %d, input %d", i, j)
		}
	}

}
