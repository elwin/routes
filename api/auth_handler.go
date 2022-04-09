package main

import (
	"net/http"
	"strings"

	app2 "github.com/elwin/heatmap/api/app"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

func (a *app) callbackHandler(c echo.Context) error {
	s, err := defaultSession(c)
	if err != nil {
		return err
	}

	code := c.QueryParam("code")
	token, err := oauthConfig(a.conf).Exchange(c.Request().Context(), code)
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
}

func (a *app) oauthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		s, err := defaultSession(c)
		if err != nil {
			return err
		}

		_, ok := s.Values["token"]
		if !ok {
			url := oauthConfig(a.conf).AuthCodeURL("state", oauth2.AccessTypeOffline)

			return c.Redirect(http.StatusSeeOther, url)
		}

		return next(c)
	}
}

func oauthConfig(conf app2.config) *oauth2.Config {
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
