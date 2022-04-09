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
		in        []Route
		out       []Route
		maxWidth  int
		maxHeight int
	}{
		{
			[]Route{
				{0, []Position{
					{0 + 50, 0 + 50},
					{10 + 50, 10 + 50},
				}},
				{0, []Position{
					{100 + 50, 90 + 50},
					{100 + 50, 100 + 50},
				}},
			},
			[]Route{
				{0, []Position{
					{0, 0},
					{1, 1},
				}},
				{0, []Position{
					{10, 9},
					{10, 10},
				}},
			},
			10,
			10,
		},
	}

	for i, test := range tests {
		assert.Equal(t, test.out, normalize(test.in, test.maxWidth, test.maxHeight), i)
	}
}
