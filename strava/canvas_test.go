package strava

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_canvas(t *testing.T) {
	activities, err := LoadActivity("tests/Activites.json")
	require.NoError(t, err)

	_, err = Draw(activities, LayoutA2)
	assert.NoError(t, err)
}
