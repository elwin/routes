package main

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/elwin/strava-go-api/v3/strava"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/pflag"
	"golang.org/x/oauth2"
)

const height, width = 1000, 4000

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	configPath := pflag.StringP("config", "c", "config.yml", "Path of config file")
	pflag.Parse()

	conf, err := readConfig(*configPath)
	if err != nil {
		return err
	}

	scopes := []string{
		"read",
		"read_all",
		"activity:read",
		"activity:read_all",
		"profile:read_all",
	}

	oauthConf := oauth2.Config{
		ClientID:     conf.Strava.ID,
		ClientSecret: conf.Strava.Secret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.strava.com/api/v3/oauth/authorize",
			TokenURL: "https://www.strava.com/api/v3/oauth/token",
		},
		RedirectURL: conf.Host + "/auth/redirect",
		Scopes:      []string{strings.Join(scopes, ",")},
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
	// 	return func(c echo.Context) error {
	//
	//
	// 		return next(c)
	// 	}
	// })

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "auth/register")
	})

	e.GET("auth/register", func(c echo.Context) error {
		url := oauthConf.AuthCodeURL("state", oauth2.AccessTypeOffline)

		return c.Redirect(http.StatusSeeOther, url)
	})

	e.GET("auth/redirect", func(c echo.Context) error {
		code := c.QueryParam("code")
		tok, err := oauthConf.Exchange(ctx, code)
		if err != nil {
			return err
		}

		conf := strava.NewConfiguration()
		conf.HTTPClient = oauthConf.Client(ctx, tok)

		stravaClient := client{strava.NewAPIClient(conf)}

		if err := stravaClient.heatMap(c.Request().Context(), width, height); err != nil {
			return err
		}

		return c.String(http.StatusOK, "all good")
	})

	// e.Use(session.Middleware(
	// 	sessions.NewFilesystemStore("/tmp/sessions", []byte("secret"))),
	// )

	// e.GET("/session", func(c echo.Context) error {
	// 	sess, err := session.Get("default", c)
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	if _, ok := sess.Values["token"]; !ok{
	// 		return c.Redirect(http.StatusSeeOther, "/auth/register")
	// 	}
	//
	// 	return nil
	// })

	return e.Start(":3030")
}

type savedToken struct {
	AccessToken  string
	Expiry       time.Time
	RefreshToken string
	TokenType    string
}
