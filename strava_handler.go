package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/elwin/heatmap/random"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func (a *app) athleteInfo(c echo.Context) error {
	client, err := newStravaClient(c, a.conf)
	if err != nil {
		return err
	}

	athlete, err := client.something(c.Request().Context())
	if err != nil {
		return err
	}

	fmt.Print(athlete)

	return c.Render(http.StatusOK, "dashboard.gohtml", "asdf")

	// return c.JSON(http.StatusOK, athlete)
}

func (a *app) imageHandler(c echo.Context) error {
	id := c.Param("id")
	width := intQuery(c, "width", defaultWidth)
	height := intQuery(c, "height", defaultHeight)
	omitted, err := intSliceQuery(c, "omit")
	if err != nil {
	    return  err
	}

	token, ok := a.db.maps[id]
	if !ok {
		return echo.ErrNotFound
	}

	if width > 10000 || height > 10000 {
		return errors.New("height and width must not exceed 10000")
	}

	client, err := newStravaClientFromToken(c.Request().Context(), token, a.conf)
	if err != nil {
		return err
	}

	image, err := client.heatMap(c.Request().Context(), width, height, omitted)
	if err != nil {
		return err
	}

	return c.Blob(http.StatusOK, "image/png", image)
}

func (a *app) temporaryImageHandler(c echo.Context) error {
	client, err := newStravaClient(c, a.conf)
	if err != nil {
		return err
	}

	width, height := 10000, 10000

	omitted, err := intSliceQuery(c, "omit")
	if err != nil {
		return  err
	}

	image, err := client.heatMap(c.Request().Context(), width, height, omitted)
	if err != nil {
		return err
	}

	return c.Blob(http.StatusOK, "image/png", image)
}

func (a *app) enableHandler(c echo.Context) error {
	client, err := newStravaClient(c, a.conf)
	if err != nil {
		return err
	}

	token := random.Generate(10)

	a.db.maps[token] = client.token

	return c.JSON(http.StatusOK, token)
}

func intQuery(c echo.Context, name string, defaultValue int) int {
	out := c.QueryParam(name)
	if out == "" {
		return defaultValue
	}

	val, err := strconv.Atoi(out)
	if err != nil {
		return defaultValue
	}

	return val
}

func intSliceQuery(c echo.Context, name string) ([]int64, error) {
	var omitted []int64

	omittedString := c.QueryParams()["omit"]
	for _, activity := range omittedString {
		i, err := strconv.ParseInt(activity, 10, 64)
		if err != nil {
			return nil, err
		}

		omitted = append(omitted, i)
	}

	return omitted, nil
}
