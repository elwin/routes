package strava

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"testing"

	"github.com/elwin/heatmap/api/app"
	"github.com/elwin/strava-go-api/v3/strava"
	"github.com/stretchr/testify/suite"
)

var (
	update = flag.Bool("update", false, "update the golden files of this test")
	port = ":3000"
	host = "http://localhost" + port
)

type StravaTestSuite struct {
	suite.Suite
	client *Client
}

type stravaConfig struct {
	Strava struct {
		ID            string `yaml:"id"`
		ClientSecret  string `yaml:"client_secret"`
		RememberID    string `yaml:"remember_id"`
		RememberToken string `yaml:"remember_token"`
	} `yaml:"strava"`
}

func (suite *StravaTestSuite) SetupTest() {
	if suite.client == nil {
		ctx := context.Background()

		var xx stravaConfig
		suite.Require().NoError(app.ReadYamlConfig("tests/strava_config.yml", &xx))

		conf := Config{
			ClientID:     xx.Strava.ID,
			ClientSecret: xx.Strava.ClientSecret,
			RedirectHost: host,
		}

		stravaRememberToken := "***REMOVED***.eyJpc3MiOiJjb20uc3RyYXZhLmF0aGxldGVzIiwic3ViIjoyMzU4NDU0MCwiaWF0IjoxNjQ3ODA2NDU3LCJleHAiOjE2NTAzOTg0NTcsImVtYWlsIjoiTTg0cExoNVFTaDF0b29vK0VpTWdhaStBV2RERUtwK1F5R0haN3YxSjc5WC9GaDNCNW9VVDJjUnJwbjVkXG5IbVpJZm8vLzE4RnNtTFhOeW8wNWM3a2hTOEJQSGtrYzJPUW8yeTdObHBlMkVFaz1cbiJ9.4cDzhb6kMi-pAXG1g5S8Dz92RcBx0NWhrdq4vSjX7S8"
		stravaRememberId := "***REMOVED***"

		token, err := FetchToken(ctx, conf, stravaRememberId, stravaRememberToken)
		suite.Require().NoError(err)

		suite.client = NewClient(ctx, conf, token)
	}
}

func (suite *StravaTestSuite) Test_activities() {
	ctx := context.Background()

	it := suite.client.activities()
	var activities []strava.SummaryActivity
	for {
		activity, err := it.Next(ctx)
		suite.Require().NoError(err)
		if activity == nil {
			break
		}

		activities = append(activities, *activity)
		if len(activities) == 5 {
			break
		}
	}

	out, err := json.MarshalIndent(activities, "", " ")
	suite.Require().NoError(err)

	if *update {
		f, err := os.Create("tests/activities.json")
		defer f.Close()
		_, err = f.Write(out)
		suite.Require().NoError(err)
	}
}

func TestStravaSuite(t *testing.T) {
	suite.Run(t, new(StravaTestSuite))
}
