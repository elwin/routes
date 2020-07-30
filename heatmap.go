package main

import (
	"github.com/elwin/strava-go-api/v3/strava"
	"github.com/fogleman/gg"
)

type client struct{
	client *strava.APIClient
}

func (c client) heatMap() error {
	activities, _, err := c.client.ActivitiesApi.GetLoggedInAthleteActivities(ctx, nil)
	if err != nil {
		return err
	}

	routes, err := convertActivitiesToRoutes(activities)
	if err != nil {
		return err
	}

	routes = filter(routes, 3316687943, 2402178038, 1095821335)

	const height, width = 750, 1334

	routes = normalize(routes, width, height)

	dc := gg.NewContext(width, height)
	dc.SetRGB(0, 0, 0)
	dc.Clear()
	dc.SetRGBA(1, 1, 1, 0.8)

	for _, route := range routes {
		if len(route.positions) == 0 {
			continue
		}

		current := route.positions[0]

		for _, next := range route.positions[1:] {
			dc.DrawLine(current.x, current.y, next.x, next.y)
			dc.Stroke()

			current = next
		}
	}

	return dc.SavePNG("out.png")
}