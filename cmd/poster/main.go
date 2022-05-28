package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/elwin/heatmap/api/strava"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
	"github.com/urfave/cli/v2"
)

func main() {
	var (
		clientID, clientSecret string
	)

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "client_id",
				Usage:       "Client ID for application",
				Required:    true,
				EnvVars:     []string{"STRAVA_CLIENT_ID"},
				Destination: &clientID,
			},
			&cli.StringFlag{
				Name:        "client_secret",
				Usage:       "Client secret for application",
				Required:    true,
				EnvVars:     []string{"STRAVA_CLIENT_SECRET"},
				Destination: &clientSecret,
			},
		},
		Name:  "poster",
		Usage: "Create a strava poster!",
		Action: func(c *cli.Context) error {
			return run(c.Context, clientID, clientSecret)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, clientID, clientSecret string) error {
	client, err := strava.New(ctx, clientID, clientSecret)
	if err != nil {
		return err
	}

	activities, err := client.Activites().All(ctx)

	l := strava.LayoutA2

	// activities, err := strava.LoadActivity("strava/tests/activities.json")
	// if err != nil {
	// 	return err
	// }

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
