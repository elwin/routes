package main

import (
	"context"
	"encoding/gob"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/elwin/strava-go-api/v3/strava"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/pflag"
	"golang.org/x/oauth2"
)

const defaultHeight, defaultWidth = 4000, 4000

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

type app struct {
	conf config
	db   *memoryDB
}

func (a *app) redirect(path string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Redirect(http.StatusSeeOther, path)
	}
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func run() error {
	configPath := pflag.StringP("config", "c", "config.yml", "Path of config file")
	pflag.Parse()

	conf, err := readConfig(*configPath)
	if err != nil {
		return err
	}

	app := app{
		conf: conf,
		db:   newMemoryDB(),
	}
	gob.Register(savedToken{})

	e := echo.New()
	e.Debug = conf.Debug
	e.Renderer = &Template{
		templates: template.Must(template.ParseGlob("resources/*.html")),
	}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(session.Middleware(sessions.NewFilesystemStore(conf.SessionDirectory, []byte("supersecret"))))
	e.Static("resources", "resources")

	e.GET("/dashboard", app.indexHandler)
	e.GET("/", app.redirect("/authorized/"))
	e.GET("/auth/redirect", app.callbackHandler)
	e.GET("/generated/:id", app.imageHandler)

	authorized := e.Group("/authorized", app.oauthMiddleware)
	authorized.GET("/", app. athleteInfo)
	authorized.GET("/image", app.temporaryImageHandler)
	authorized.GET("/enable", app.enableHandler)

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

	return newStravaClientFromToken(c.Request().Context(), token, conf)
}

func newStravaClientFromToken(ctx context.Context, token savedToken, conf config) (client, error) {
	clientConfig := strava.NewConfiguration()
	clientConfig.HTTPClient = oauthConfig(conf).Client(ctx, &oauth2.Token{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	})

	return client{
		client: strava.NewAPIClient(clientConfig),
		token:  token,
	}, nil
}
