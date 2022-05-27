package strava

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/elwin/strava-go-api/v3/strava"
	"github.com/go-resty/resty/v2"
	"golang.org/x/oauth2"
)

type Client struct {
	client *strava.APIClient
}

type Config struct {
	ClientID     string
	ClientSecret string
	RedirectHost string
}

func oauthConfig(conf Config) *oauth2.Config {
	scopes := []string{
		"read",
		"read_all",
		"activity:read",
		"activity:read_all",
		"profile:read_all",
	}

	return &oauth2.Config{
		ClientID:     conf.ClientID,
		ClientSecret: conf.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.strava.com/api/v3/oauth/authorize",
			TokenURL: "https://www.strava.com/api/v3/oauth/token",
		},
		RedirectURL: conf.RedirectHost + "/auth/redirect",
		Scopes:      []string{strings.Join(scopes, ",")},
	}
}

func FetchToken(ctx context.Context, conf Config, stravaRememberId, stravaRememberToken string) (*oauth2.Token, error) {
	url := oauthConfig(conf).AuthCodeURL("state", oauth2.AccessTypeOffline)

	handler := http.NewServeMux()
	srv := http.Server{Addr: strings.TrimLeft(conf.RedirectHost, "http://"), Handler: handler}
	var code string
	wg := &sync.WaitGroup{}
	wg.Add(1)
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code = r.URL.Query().Get("code")
		if _, err := w.Write([]byte("Authorized")); err != nil {
			log.Fatal(err)
		}

		wg.Done()
	})

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()

	resp, err := resty.New().SetCookie(&http.Cookie{
		Name:  "strava_remember_token",
		Value: stravaRememberToken,
	}).SetCookie(&http.Cookie{
		Name:  "strava_remember_id",
		Value: stravaRememberId,
	}).R().Get(url)
	if err != nil {
		return nil, err
	}

	if resp.String() != "Authorized" {
		fmt.Println("Please open the following URL in your browser")
		fmt.Println(url)
	}

	wg.Wait()

	if err := srv.Shutdown(ctx); err != nil {
		return nil, err
	}

	return oauthConfig(conf).Exchange(ctx, code)
}
func NewClient(ctx context.Context, conf Config, token *oauth2.Token) *Client {
	clientConfig := strava.NewConfiguration()
	clientConfig.HTTPClient = oauthConfig(conf).Client(ctx, token)

	return &Client{
		client: strava.NewAPIClient(clientConfig),
	}
}
