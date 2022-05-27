package strava

import (
	"embed"
	"encoding/json"
	"fmt"
	"image/color"
	"os"

	"github.com/elwin/strava-go-api/v3/strava"
	"github.com/goodsign/monday"
	"github.com/icza/gox/imagex/colorx"
	"github.com/samber/lo"
	"github.com/tdewolff/canvas"
	"go.uber.org/zap"
)

var (
	//go:embed resources/fonts
	fonts embed.FS

	neon    = ColorPalette{"neon", "#FFFFFF", "#541690", "#FF8D29", "#FF4949", "#C6CEDB"}
	vintage = ColorPalette{"vintage", "#F7FAFF", "#DE8971", "#A7D0CD", "#7B6079", "#FFFFFF"}
	bright  = ColorPalette{"bright", "#F7FAFF", "#9ADCFF", "#94B49F", "#FFB2A6", "#FFFFFF"}
	dawn    = ColorPalette{"dawn", "#180A0A", "#711A75", "#F10086", "#F582A7", "#FFFFFF"}
	sea     = ColorPalette{"sea", "#F7E2E2", "#61A4BC", "#5B7DB1", "#1A132F", "#FFFFFF"}
	happy   = ColorPalette{"happy", "#FF6B6B", "#FFD93D", "#6BCB77", "#4D96FF", "#FFFFFF"}
	hotdog  = ColorPalette{"hotdog", "#333C83", "#F24A72", "#FDAF75", "#EAEA7F", "#FFFFFF"}
	dark    = ColorPalette{"dark", "#1A1A2E", "#16213E", "#0F3460", "#E94560", "#FFFFFF"}
	violet  = ColorPalette{"violet", "#000000", "#52057B", "#892CDC", "#BC6FF1", "#FFFFFF"}
	red     = ColorPalette{"red", "#082032", "#2C394B", "#334756", "#FF4C29", "#FFFFFF"}
	blood   = ColorPalette{"blood", "#000000", "#3D0000", "#950101", "#FF0000", "#FFFFFF"}
	aga     = ColorPalette{"aga", "#2F4858", "#F26419", "#55DDE0", "#F6AE2D", "#FFFFFF"}

	palettes = []ColorPalette{neon, vintage, bright, dawn, sea, happy, hotdog, dark, violet, red, blood, aga}
	LayoutA2 = Layout{
		HorizontalBoxes:  5,
		VerticalBoxes:    9,
		BoxHeight:        400,
		BoxWidth:         400,
		HorizontalMargin: 400,
		VerticalMargin:   800,
		TotalWidth:       4200,
		TotalHeight:      5940,
		ColorPalette:     neon,
		Debug:            false,
		PrintText:        true,
		Fontsize:         64,
	}

	log *zap.Logger
)

func init() {
	var err error
	log, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
}

type Layout struct {
	HorizontalBoxes, VerticalBoxes   int
	BoxHeight, BoxWidth              float64
	HorizontalMargin, VerticalMargin float64
	TotalWidth, TotalHeight          float64
	ColorPalette                     ColorPalette
	Debug, PrintText                 bool
	Fontsize                         float64
}

func (l Layout) horizontalGap() float64 {
	return (l.TotalWidth - 2*l.HorizontalMargin - float64(l.HorizontalBoxes)*l.BoxWidth) / (float64(l.HorizontalBoxes) - 1)
}

func (l Layout) verticalGap() float64 {
	return (l.TotalHeight - 2*l.VerticalMargin - float64(l.VerticalBoxes)*l.BoxHeight) / (float64(l.VerticalBoxes) - 1)
}

func parseHexColor(s string) color.RGBA {
	c, err := colorx.ParseHexColor(s)
	if err != nil {
		panic(err)
	}

	return c
}

type ColorPalette struct {
	name, background, ride, winter, run, font string
}

func (p ColorPalette) String() string {
	return p.name
}

func (p ColorPalette) generateColorMap() map[strava.ActivityType]color.RGBA {
	return map[strava.ActivityType]color.RGBA{
		strava.RIDE_ActivityType:            parseHexColor(p.ride),
		strava.ALPINE_SKI_ActivityType:      parseHexColor(p.winter),
		strava.BACKCOUNTRY_SKI_ActivityType: parseHexColor(p.winter),
		strava.NORDIC_SKI_ActivityType:      parseHexColor(p.winter),
		strava.RUN_ActivityType:             parseHexColor(p.run),
		strava.WALK_ActivityType:            parseHexColor(p.run),
		strava.SNOWSHOE_ActivityType:        parseHexColor(p.run),
		strava.HIKE_ActivityType:            parseHexColor(p.run),
	}
}

func (p ColorPalette) backgroundColor() color.RGBA {
	return parseHexColor(p.background)
}

func (p ColorPalette) fontColor() color.RGBA {
	return parseHexColor(p.font)
}

func LoadActivity(path string) ([]strava.SummaryActivity, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var activities []strava.SummaryActivity
	return activities, json.NewDecoder(f).Decode(&activities)
}

func Draw(activities []strava.SummaryActivity, layout Layout) (*canvas.Canvas, error) {
	colorMap := layout.ColorPalette.generateColorMap()

	c := canvas.New(layout.TotalWidth, layout.TotalHeight)
	ctx := canvas.NewContext(c)
	ctx.SetFillColor(layout.ColorPalette.backgroundColor())
	ctx.DrawPath(0, 0, canvas.Rectangle(layout.TotalWidth, layout.TotalHeight))

	blacklist := map[int64]bool{
		6846376734: true,
		5646108624: true,
		5716991890: true,
		6635240183: true,
	}

	activities = lo.Filter[strava.SummaryActivity](activities, func(a strava.SummaryActivity, _ int) bool {
		return a.Map_.SummaryPolyline != "" && !a.Private && !blacklist[a.Id]
	})
	items := lo.Min([]int{len(activities), layout.VerticalBoxes * layout.HorizontalBoxes})
	activities = lo.Reverse(activities[:items])

	family := canvas.NewFontFamily("sf-compact")
	requiredFonts := map[string]canvas.FontStyle{
		"SourceSansPro-Regular.ttf":    canvas.FontRegular,
		"SourceSansPro-ExtraLight.ttf": canvas.FontExtraLight,
	}

	for path, style := range requiredFonts {
		f, err := fonts.ReadFile(fmt.Sprintf("resources/fonts/%s", path))
		if err != nil {
			return nil, err
		}

		if err := family.LoadFont(f, 0, style); err != nil {
			return nil, err
		}
	}

	faceRegular := family.Face(layout.Fontsize, layout.ColorPalette.fontColor(), canvas.FontRegular, canvas.FontNormal)
	faceLight := family.Face(layout.Fontsize, layout.ColorPalette.fontColor(), canvas.FontExtraLight, canvas.FontNormal)

	for i := 0; i < layout.VerticalBoxes; i++ {
		for j := 0; j < layout.HorizontalBoxes; j++ {
			idx := i*layout.HorizontalBoxes + j
			if idx >= items {
				break
			}

			x := layout.HorizontalMargin + (layout.BoxWidth+layout.horizontalGap())*float64(j)
			y := layout.TotalHeight - layout.BoxHeight - (layout.VerticalMargin + (layout.BoxHeight+layout.verticalGap())*float64(i))

			activity := activities[idx]

			if layout.Debug {
				ctx.SetStrokeWidth(1)
				ctx.SetStrokeColor(canvas.Red)
				ctx.DrawPath(x, y, canvas.Rectangle(layout.BoxWidth, layout.BoxHeight))
			}

			if layout.PrintText {
				date := monday.Format(activity.StartDate, "2. January 2006", "de_DE")
				activityName := removeEmojis(activity.Name)
				ctx.DrawText(x, y-20, canvas.NewTextBox(faceLight, date, layout.BoxWidth, 0, canvas.Top, canvas.Left, 0, 0))
				ctx.DrawText(x, y-48, canvas.NewTextBox(faceRegular, activityName, layout.BoxWidth, 0, canvas.Top, canvas.Left, 0, 0))
			}

			route, err := convertActivityToRoute(activity)
			if err != nil {
				return nil, err
			}

			route = normalize(route, layout.BoxWidth, layout.BoxHeight)

			path := &canvas.Path{}
			path.MoveTo(route.Positions[0].X, route.Positions[0].Y)
			for _, p := range route.Positions[1:] {
				path.LineTo(p.X, p.Y)
			}

			var curColor color.RGBA
			if typeColor, ok := colorMap[*activity.Type_]; ok {
				curColor = typeColor
			} else {
				curColor = canvas.Red
				log.Warn("Activity w/o explicit color", zap.String("activity", string(*(activity.Type_))))
			}

			ctx.SetStrokeWidth(4)
			ctx.SetStrokeColor(curColor)
			ctx.DrawPath(x, y, path)
		}
	}

	return c, nil
}
