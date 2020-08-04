package main

import (
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

	return c.JSON(http.StatusOK, athlete)
}

func (a *app) imageHandler(c echo.Context) error {
	id := c.Param("id")
	fileType := c.QueryParam("type")
	width := intQuery(c, "width", defaultWidth)
	height := intQuery(c, "height", defaultHeight)

	token, ok := a.db.maps[id]
	if !ok {
		return echo.ErrNotFound
	}

	if width > 10000 || height > 10000 {
		return errors.New("height and width must not exceed 10000")
	}

	if fileType != "png" {
		return echo.ErrNotFound
	}
	client, err := newStravaClientFromToken(c.Request().Context(), token, a.conf)
	if err != nil {
		return err
	}

	image, err := client.heatMap(c.Request().Context(), width, height)
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

	width, height := 1000, 1000

	image, err := client.heatMap(c.Request().Context(), width, height)
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
