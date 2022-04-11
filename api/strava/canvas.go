package strava

import (
	"encoding/json"
	"fmt"
	"image/color"
	"os"

	"github.com/elwin/strava-go-api/v3/strava"
	"github.com/icza/gox/imagex/colorx"
	"github.com/samber/lo"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
)

const (
	horizontalBoxes = 5
	verticalBoxes   = 10

	boxHeight = 200
	boxWidth  = 200

	horizontalMargin = 200
	verticalMargin   = 200

	totalWidth  = 2000
	totalHeight = 3000

	horizontalGap = (totalWidth - 2*horizontalMargin - horizontalBoxes*boxWidth) / (horizontalBoxes - 1)
	verticalGap   = (totalHeight - 2*verticalMargin - verticalBoxes*boxHeight) / (verticalBoxes - 1)

	debug = false
)

func parseHexColor(s string) color.RGBA {
	c, err := colorx.ParseHexColor(s)
	if err != nil {
		panic(err)
	}

	return c
}

type colorPalette struct {
	name, background, ride, ski, run string
}

var (
	backgroundColor = parseHexColor("#011C27")
	colorMap        = map[strava.ActivityType]color.RGBA{
		strava.RIDE_ActivityType:            parseHexColor("#03254E"),
		strava.ALPINE_SKI_ActivityType:      parseHexColor("#545677"),
		strava.BACKCOUNTRY_SKI_ActivityType: parseHexColor("#545677"),
		strava.NORDIC_SKI_ActivityType:      parseHexColor("#545677"),
		strava.RUN_ActivityType:             parseHexColor("#EB9FEF"),
	}

	palettes = []colorPalette{
		{"dawn", "#180A0A", "#711A75", "#F10086", "#F582A7"},
		{"sea", "#F7E2E2", "#61A4BC", "#5B7DB1", "#1A132F"},
		{"happy", "#FF6B6B", "#FFD93D", "#6BCB77", "#4D96FF"},
		{"hotdog", "#333C83", "#F24A72", "#FDAF75", "#EAEA7F"},
		{"dark", "#1A1A2E", "#16213E", "#0F3460", "#E94560"},
		{"violet", "#000000", "#52057B", "#892CDC", "#BC6FF1"},
		{"strava", "#082032", "#2C394B", "#334756", "#FF4C29"},
		{"blood", "#000000", "#3D0000", "#950101", "#FF0000"},
	}
)

func load() ([]strava.SummaryActivity, error) {
	f, err := os.Open("tests/activities.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var activities []strava.SummaryActivity
	return activities, json.NewDecoder(f).Decode(&activities)
}

func draw(activities []strava.SummaryActivity, palette colorPalette) error {
	colorMap = map[strava.ActivityType]color.RGBA{
		strava.RIDE_ActivityType:            parseHexColor(palette.ride),
		strava.ALPINE_SKI_ActivityType:      parseHexColor(palette.ski),
		strava.BACKCOUNTRY_SKI_ActivityType: parseHexColor(palette.ski),
		strava.NORDIC_SKI_ActivityType:      parseHexColor(palette.ski),
		strava.RUN_ActivityType:             parseHexColor(palette.run),
		strava.WALK_ActivityType:            parseHexColor(palette.run),
		strava.SNOWSHOE_ActivityType:        parseHexColor(palette.run),
	}

	c := canvas.New(totalWidth, totalHeight)
	ctx := canvas.NewContext(c)
	ctx.SetFillColor(parseHexColor(palette.background))
	ctx.DrawPath(0, 0, canvas.Rectangle(totalWidth, totalHeight))

	if debug {
		ctx.SetFillColor(canvas.Transparent)
		ctx.SetStrokeWidth(1)

		for i := 0; i < horizontalBoxes; i++ {
			for j := 0; j < verticalBoxes; j++ {
				x := float64(300 + 500*i)
				y := float64(300 + 500*j)

				ctx.DrawPath(x, y, canvas.Rectangle(boxWidth, boxHeight))
			}
		}
	}

	activities = lo.Filter[strava.SummaryActivity](activities, func(a strava.SummaryActivity, _ int) bool {
		return a.Map_.SummaryPolyline != "" && !a.Private
	})

	ctx.SetStrokeWidth(4)
	ctx.SetStrokeColor(canvas.Red)

	for i := 0; i < verticalBoxes; i++ {
		for j := 0; j < horizontalBoxes; j++ {
			idx := i*verticalBoxes + j
			if idx >= len(activities) {
				break
			}

			activity := activities[idx]
			x := float64(horizontalMargin + (boxWidth+horizontalGap)*j)
			y := totalHeight - boxHeight - float64(verticalMargin+(boxHeight+verticalGap)*i)

			route, err := convertActivityToRoute(activity)
			if err != nil {
				return err
			}

			route = norm(route, boxWidth, boxHeight)

			path := &canvas.Path{}
			path.MoveTo(route.Positions[0].X, route.Positions[0].Y)
			for _, p := range route.Positions[1:] {
				path.LineTo(p.X, p.Y)
			}

			curColor := canvas.White
			if typeColor, ok := colorMap[*activity.Type_]; ok {
				curColor = typeColor
			} else {
				fmt.Println(*(activity.Type_))
			}

			ctx.SetStrokeColor(curColor)
			ctx.DrawPath(x, y, path)
		}
	}

	filename := palette.name
	if err := c.WriteFile("out/png/"+filename+".png", renderers.PNG()); err != nil {
		return err
	}
	if err := c.WriteFile("out/pdf/"+filename+".pdf", renderers.PDF()); err != nil {
		return err
	}

	return nil
}
