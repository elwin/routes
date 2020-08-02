package main

import (
	"context"
	"encoding/gob"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/elwin/strava-go-api/v3/strava"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/pflag"
	"golang.org/x/oauth2"
)

const height, width = 4000, 4000

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

	e := echo.New()
	e.Debug = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(session.Middleware(sessions.NewFilesystemStore("/tmp/sessions", []byte("supersecret"))))

	authorized := e.Group("/authorized", func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			s, err := defaultSession(c)
			if err != nil {
				return err
			}

			_, ok := s.Values["token"]
			if !ok {
				url := oauthConfig(conf).AuthCodeURL("state", oauth2.AccessTypeOffline)

				return c.Redirect(http.StatusSeeOther, url)
			}

			return next(c)
		}
	})

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/authorized/")
	})

	e.GET("/auth/redirect", func(c echo.Context) error {
		s, err := defaultSession(c)
		if err != nil {
			return err
		}

		code := c.QueryParam("code")
		token, err := oauthConfig(conf).Exchange(ctx, code)
		if err != nil {
			return err
		}

		s.Values["token"] = savedToken{
			AccessToken:  token.AccessToken,
			Expiry:       token.Expiry,
			RefreshToken: token.RefreshToken,
			TokenType:    token.TokenType,
		}
		if err := s.Store().Save(c.Request(), c.Response(), s); err != nil {
			return err
		}

		return c.Redirect(http.StatusSeeOther, "/authorized/")
	})

	authorized.GET("/", func(c echo.Context) error {
		client, err := newStravaClient(c, conf)
		if err != nil {
			return err
		}

		athlete, err := client.something(c.Request().Context())
		if err != nil {
		    return  err
		}

		return c.JSON(http.StatusOK, athlete)
	})

	authorized.GET("/image", func(c echo.Context) error {
		client, err := newStravaClient(c, conf)
		if err != nil {
			return err
		}

		img, err := client.heatMap(c.Request().Context(), width, height)
		if err != nil {
			return err
		}

		return c.Blob(http.StatusOK, "image/png", img)
	})

	gob.Register(savedToken{})

	return e.Start(":3030")
}

type savedToken struct {
	AccessToken  string
	Expiry       time.Time
	RefreshToken string
	TokenType    string
}

func defaultSession(c echo.Context) (*sessions.Session, error) {
	return session.Get("default", c)
}

func newStravaClient(c echo.Context, conf config) (client, error) {
	s, err := defaultSession(c)
	if err != nil {
		return client{}, err
	}

	token := s.Values["token"].(savedToken)

	clientConfig := strava.NewConfiguration()
	clientConfig.HTTPClient = oauthConfig(conf).Client(c.Request().Context(), &oauth2.Token{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	})

	return client{strava.NewAPIClient(clientConfig)}, nil
}

func oauthConfig(conf config) *oauth2.Config {
	scopes := []string{
		"read",
		"read_all",
		"activity:read",
		"activity:read_all",
		"profile:read_all",
	}

	return &oauth2.Config{
		ClientID:     conf.Strava.ID,
		ClientSecret: conf.Strava.Secret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.strava.com/api/v3/oauth/authorize",
			TokenURL: "https://www.strava.com/api/v3/oauth/token",
		},
		RedirectURL: conf.Host + "/auth/redirect",
		Scopes:      []string{strings.Join(scopes, ",")},
	}
}
