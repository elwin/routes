package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (a *app) indexHandler(c echo.Context) error {

	return c.Render(http.StatusOK, "index.html", nil)

	// return c.JSON(200, "yo")

	// return nil
}
