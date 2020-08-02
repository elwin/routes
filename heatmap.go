package main

import (
	"bytes"
	"context"

	"github.com/antihax/optional"
	"github.com/elwin/strava-go-api/v3/strava"
	"github.com/fogleman/gg"
	"github.com/pkg/errors"
	"github.com/twpayne/go-polyline"
)

type client struct {
	client *strava.APIClient
	token  savedToken
}

func (c client) activities(ctx context.Context) ([]strava.SummaryActivity, error) {
	currentPage := 1

	var result []strava.SummaryActivity

	for {
		options := &strava.ActivitiesApiGetLoggedInAthleteActivitiesOpts{
			Page: optional.NewInt32(int32(currentPage)),
		}
		activities, _, err := c.client.ActivitiesApi.GetLoggedInAthleteActivities(ctx, options)
		if err != nil {
			return nil, err
		}

		if len(activities) == 0 {
			break
		}

		result = append(result, activities...)
		currentPage++
	}

	return result, nil
}

func (c client) heatMap(ctx context.Context, maxWidth, maxHeight int) ([]byte, error) {
	activities, err := c.activities(ctx)
	if err != nil {
	    return nil, err
	}

	routes, err := convertActivitiesToRoutes(activities)
	if err != nil {
		return nil, err
	}

	routes = filter(routes, 3316687943, 2402178038, 1095821335)

	routes = normalize(routes, maxWidth, maxHeight)

	dc := gg.NewContext(maxWidth, maxHeight)
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

	buf := bytes.NewBuffer(nil)
	if err := dc.EncodePNG(buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func convertActivitiesToRoutes(activities []strava.SummaryActivity) ([]route, error) {
	var routes []route
	for _, activity := range activities {

		// skip activities without a tracked route
		if activity.Map_.SummaryPolyline == "" {
			continue
		}

		coords, _, err := polyline.DecodeCoords([]byte(activity.Map_.SummaryPolyline))
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode coordinates")
		}

		var positions []position
		for _, coord := range coords {
			if len(coord) != 2 {
				return nil, errors.Errorf("expected 2 coordinates (x, y), received %d", len(coord))
			}

			positions = append(positions, position{x: coord[1], y: -coord[0]})
		}

		routes = append(routes, route{
			id:        activity.Id,
			positions: positions,
		})
	}

	return routes, nil
}

func (c client) something(ctx context.Context) (strava.ActivityStats, error) {
	athlete, _, err := c.client.AthletesApi.GetLoggedInAthlete(ctx)
	if err != nil {
		return strava.ActivityStats{}, err
	}

	stats, _, err := c.client.AthletesApi.GetStats(ctx, athlete.Id)
	if err != nil {
		return strava.ActivityStats{}, err
	}

	return stats, nil

}
