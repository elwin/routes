package main

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/elwin/routes/strava"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
)

var (
	clientID, clientSecret string
	// rememberID, rememberToken string
	outputPath string

	stravaPull = &cli.Command{
		Name:        "strava",
		Description: "Download all activities from Strava",
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
			// &cli.StringFlag{
			// 	Name:        "remember_id",
			// 	Destination: &rememberID,
			// 	EnvVars:     []string{"STRAVA_REMEMBER_ID"},
			// },
			// &cli.StringFlag{
			// 	Name:        "remember_token",
			// 	Destination: &rememberToken,
			// 	EnvVars:     []string{"STRAVA_REMEMBER_TOKEN"},
			// },
			&cli.StringFlag{
				Name:        "output",
				Aliases:     []string{"o"},
				Destination: &outputPath,
			},
		},
		Action: func(c *cli.Context) error {
			client, err := strava.New(c.Context, clientID, clientSecret)
			if err != nil {
				return err
			}

			activities, err := client.Activites().All(c.Context)

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Name", "Start Date"})
			for _, activity := range activities {
				table.Append([]string{
					strconv.FormatInt(activity.Id, 10),
					activity.Name,
					activity.StartDate.Format("2. January 2006"),
				})
			}
			table.Render()

			f, err := os.Create(outputPath)
			defer f.Close()

			return json.NewEncoder(f).Encode(activities)
		},
	}
)
