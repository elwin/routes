package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/elwin/routes/strava"
	stravaApi "github.com/elwin/strava-go-api/v3/strava"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
	"github.com/urfave/cli/v2"
)

var (
	path string

	poster = &cli.Command{
		Name:        "poster",
		Description: "Convert Strava activities to a poster",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "input",
				Usage:       "Path for activities.json file",
				Required:    true,
				Destination: &path,
			},
		},
		Action: func(ctx *cli.Context) error {
			f, err := os.Open(path)
			if err != nil {
				return err
			}

			var activities []stravaApi.SummaryActivity
			if err := json.NewDecoder(f).Decode(&activities); err != nil {
				return err
			}

			c, err := strava.Draw(activities, strava.LayoutA2)
			if err != nil {
				return err
			}

			filename := "poster"
			outputs := map[string]canvas.Writer{
				"png": renderers.PNG(),
				"pdf": renderers.PDF(),
				"svg": renderers.SVG(),
			}

			for output, renderer := range outputs {
				if err := os.MkdirAll(fmt.Sprintf("%s", output), 0755); err != nil {
					return err
				}

				if err := c.WriteFile(fmt.Sprintf("%s.%s", filename, output), renderer); err != nil {
					return err
				}
			}

			return nil
		},
	}
)
