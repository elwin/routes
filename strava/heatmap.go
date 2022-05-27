package strava

import (
	"context"

	"github.com/antihax/optional"
	"github.com/elwin/strava-go-api/v3/strava"
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
