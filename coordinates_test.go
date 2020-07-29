package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_bounds(t *testing.T) {
	tests := []struct {
		in          []position
		topLeft     position
		bottomRight position
	}{
		{
			[]position{
				{1, 1},
				{4, 9},
				{6, 3},
				{9, 3},
			},
			position{1, 1},
			position{9, 9},
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
		in        []route
		out       []route
		maxWidth  float64
		maxHeight float64
	}{
		{
			[]route{
				{[]position{
					{0 + 50, 0 + 50},
					{10 + 50, 10 + 50},
				}},
				{[]position{
					{100 + 50, 90 + 50},
					{100 + 50, 100 + 50},
				}},
			},
			[]route{
				{[]position{
					{0, 0},
					{1, 1},
				}},
				{[]position{
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
