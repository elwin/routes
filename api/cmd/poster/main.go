package main

import (
	"fmt"
	"log"
	"os"

	"github.com/elwin/heatmap/api/strava"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	l := strava.LayoutA2

	activities, err := strava.LoadActivity("strava/tests/activities.json")
	if err != nil {
		return err
	}

	c, err := strava.Draw(activities, l)
	if err != nil {
		return err
	}

	filename := l.ColorPalette.String()
	outputs := map[string]canvas.Writer{
		"png": renderers.PNG(),
		"pdf": renderers.PDF(),
		"svg": renderers.SVG(),
	}

	for output, renderer := range outputs {
		if err := os.MkdirAll(fmt.Sprintf("out/%s", output), 0755); err != nil {
			return err
		}

		if err := c.WriteFile(fmt.Sprintf("out/%s/%s.%s", output, filename, output), renderer); err != nil {
			return err
		}
	}

	return nil
}
