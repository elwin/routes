package strava

import (
	"bytes"
	"context"

	"github.com/antihax/optional"
	"github.com/elwin/strava-go-api/v3/strava"
	"github.com/fogleman/gg"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/twpayne/go-polyline"
)

type ActivityIterator struct {
	c           Client
	currentPage int
	buffer      []strava.SummaryActivity
}

func (it *ActivityIterator) Next(ctx context.Context) (*strava.SummaryActivity, error) {
	if len(it.buffer) == 0 {
		options := &strava.ActivitiesApiGetLoggedInAthleteActivitiesOpts{
			Page: optional.NewInt32(int32(it.currentPage)),
		}
		activities, _, err := it.c.client.ActivitiesApi.GetLoggedInAthleteActivities(ctx, options)
		if err != nil {
			return nil, err
		}

		if len(activities) == 0 {
			return nil, nil
		}

		it.buffer = append(it.buffer, lo.Reverse[strava.SummaryActivity](activities)...)
		it.currentPage++
	}

	next := it.buffer[len(it.buffer)-1]
	it.buffer = it.buffer[:len(it.buffer)-1]

	return &next, nil
}

func (it *ActivityIterator) all(ctx context.Context) ([]strava.SummaryActivity, error) {
	var results []strava.SummaryActivity
	for {
		if activity, err := it.Next(ctx); err != nil {
			return nil, err
		} else if activity == nil {
			return results, nil
		} else {
			results = append(results, *activity)
		}
	}
}

func (c Client) activities() *ActivityIterator {
	return &ActivityIterator{
		c:           c,
		currentPage: 1,
		buffer:      []strava.SummaryActivity{},
	}
}

func (c Client) heatmap(ctx context.Context, maxWidth, maxHeight int, omit []int64) ([]byte, error) {
	activities, err := c.activities().all(ctx)
	if err != nil {
		return nil, err
	}

	routes, err := convertActivitiesToRoutes(activities)
	if err != nil {
		return nil, err
	}

	routes = filter(routes, omit...)

	routes = normalize(routes, maxWidth, maxHeight)

	dc := gg.NewContext(maxWidth, maxHeight)
	dc.SetRGB(0, 0, 0)
	dc.Clear()
	dc.SetRGBA(1, 1, 1, 0.8)

	for _, route := range routes {
		if len(route.Positions) == 0 {
			continue
		}

		current := route.Positions[0]

		for _, next := range route.Positions[1:] {
			dc.DrawLine(current.X, current.Y, next.X, next.Y)
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

func convertActivityToRoute(activity strava.SummaryActivity) (Route, error) {
	coords, _, err := polyline.DecodeCoords([]byte(activity.Map_.SummaryPolyline))
	if err != nil {
		return Route{}, errors.Wrap(err, "failed to decode coordinates")
	}

	var positions []Position
	for _, coord := range coords {
		if len(coord) != 2 {
			return Route{}, errors.Errorf("expected 2 coordinates (X, Y), received %d", len(coord))
		}

		positions = append(positions, Position{X: coord[1], Y: -coord[0]})
	}

	return Route{
		Id:        activity.Id,
		Positions: positions,
	}, nil
}

func convertActivitiesToRoutes(activities []strava.SummaryActivity) ([]Route, error) {
	var routes []Route
	for _, activity := range activities {

		// skip activities without a tracked Route
		if activity.Map_.SummaryPolyline == "" {
			continue
		}

		coords, _, err := polyline.DecodeCoords([]byte(activity.Map_.SummaryPolyline))
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode coordinates")
		}

		var positions []Position
		for _, coord := range coords {
			if len(coord) != 2 {
				return nil, errors.Errorf("expected 2 coordinates (X, Y), received %d", len(coord))
			}

			positions = append(positions, Position{X: coord[1], Y: -coord[0]})
		}

		routes = append(routes, Route{
			Id:        activity.Id,
			Positions: positions,
		})
	}

	return routes, nil
}

// func (c StravaClient) something(ctx context.Context) (strava.ActivityStats, error) {
// 	athlete, _, err := c.StravaClient.AthletesApi.GetLoggedInAthlete(ctx)
// 	if err != nil {
// 		return strava.ActivityStats{}, err
// 	}
//
// 	stats, _, err := c.StravaClient.AthletesApi.GetStats(ctx, athlete.Id)
// 	if err != nil {
// 		return strava.ActivityStats{}, err
// 	}
//
// 	return stats, nil
//
// }
