package strava

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_canvas(t *testing.T) {
	activities, err := load()
	assert.NoError(t, err)

	assert.NoError(t, os.MkdirAll("out/png", 0755))
	assert.NoError(t, os.MkdirAll("out/pdf", 0755))
	assert.NoError(t, os.MkdirAll("out/svg", 0755))


	for _, palette := range palettes {
		assert.NoError(t, draw(activities, palette))
	}
}
