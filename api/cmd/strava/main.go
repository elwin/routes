package main

import (
	"context"
	"fmt"
	"log"

	"github.com/elwin/heatmap/api/strava"
	"github.com/stretchr/testify/assert"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	conf := strava.Config{
		ClientID: "***REMOVED***",
		ClientSecret: "***REMOVED***",
		RedirectHost: "http://localhost:8080",
	}

	token, err := strava.FetchToken(ctx, conf)
	if err != nil {
	    return  err
	}

	client := strava.NewClient(ctx, conf, token)

	it := client.Routes()
	for {
		activity, err := it.next(ctx)
		assert.NoError(t, err)
		if activity == nil {
			break
		}

		fmt.Println(activity.AchievementCount)
	}

	assert.NoError(t, err)
}