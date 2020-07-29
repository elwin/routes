package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/elwin/strava-go-api/v3/strava"
	"github.com/fogleman/gg"
	"github.com/twpayne/go-polyline"
	"golang.org/x/oauth2"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	scopes := []string{
		"read",
		"read_all",
		"activity:read",
		"activity:read_all",
		"profile:read_all",
	}

	oauthConf := oauth2.Config{
		ClientID:     "***REMOVED***",
		ClientSecret: "e13bd1b6799d074e354443129091ffe40e2add96",
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://www.strava.com/api/v3/oauth/authorize",
			TokenURL:  "https://www.strava.com/api/v3/oauth/token",
			AuthStyle: 0,
		},
		RedirectURL: "http://localhost:3030/return",
		Scopes:      []string{strings.Join(scopes, ",")},
	}

	register := func() http.Handler {
		url := oauthConf.AuthCodeURL("state", oauth2.AccessTypeOffline)

		return http.RedirectHandler(url, http.StatusSeeOther)
	}

	h := func(f func(r *http.Request) (string, error)) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			success, err := f(r)
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}

			w.Write([]byte(success))
		}
	}

	http.Handle("/", register())
	http.HandleFunc("/return", h(func(r *http.Request) (s string, err error) {
		code := r.URL.Query().Get("code")
		if code == "" {

		}

		tok, err := oauthConf.Exchange(ctx, code)
		if err != nil {
			return "", err
		}

		conf := strava.NewConfiguration()
		conf.HTTPClient = oauthConf.Client(ctx, tok)
		client := strava.NewAPIClient(conf)

		activites, _, err := client.ActivitiesApi.GetLoggedInAthleteActivities(ctx, nil)
		if err != nil {
			return "", err
		}

		routes, err := convertActivitiesToRoutes(activites)
		if err != nil {
			return "", err
		}

		const height, width = 4096, 4096

		routes = normalize(routes, width, height)

		var out string

		dc := gg.NewContext(width, height)
		dc.SetRGB(0, 0, 0)
		dc.Clear()
		dc.SetRGB(1, 1, 1)

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

		dc.SavePNG("out.png")

		return out, nil
	}))

	return http.ListenAndServe(":3030", http.DefaultServeMux)
}

func convertActivitiesToRoutes(activities []strava.SummaryActivity) ([]route, error) {
	var routes []route
	for _, activity := range activities {
		coords, _, err := polyline.DecodeCoords([]byte(activity.Map_.SummaryPolyline))
		if err != nil {
			return nil, err
		}

		var positions []position
		for _, coord := range coords {
			if len(coord) != 2 {
				return nil, fmt.Errorf("expected 2 coordinates (x, y), received %d", len(coord))
			}

			positions = append(positions, position{x: coord[1], y: -coord[0]})
		}

		routes = append(routes, route{positions: positions})
	}

	return routes, nil
}