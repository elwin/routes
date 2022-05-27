package strava

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_bounds(t *testing.T) {
	tests := []struct {
		in          []Position
		topLeft     Position
		bottomRight Position
	}{
		{
			[]Position{
				{1, 1},
				{4, 9},
				{6, 3},
				{9, 3},
			},
			Position{1, 1},
			Position{9, 9},
		},
	}

	for i, test := range tests {
		topLeft, bottomRight := bounds(test.in)
		assert.Equal(t, test.topLeft, topLeft, i)
		assert.Equal(t, test.bottomRight, bottomRight, i)
	}
}

func Test_normalize(t *testing.T) {
	tests := []struct {
		in, out             Route
		maxWidth, maxHeight float64
	}{
		{
			Route{0, []Position{
				{0, 0},
				{100, 100},
			}},
			Route{0, []Position{
				{0, 0},
				{10, 10},
			}},
			10,
			10,
		},
		{
			Route{0, []Position{
				{0, 0},
				{50, 50},
				{100, 100},
			}},
			Route{0, []Position{
				{0, 0},
				{5, 5},
				{10, 10},
			}},
			10,
			10,
		},
		{
			Route{
				0, []Position{
					{50 + 100, 50 + 90},
					{50 + 90, 50 + 100},
				},
			},
			Route{0, []Position{
				{10, 0},
				{0, 10},
			}},
			10,
			10,
		},
	}

	for i, test := range tests {
		assert.Equal(t, test.out, normalize(test.in, test.maxWidth, test.maxHeight), i)
	}
}
