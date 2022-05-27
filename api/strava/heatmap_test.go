package strava

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"testing"

	"github.com/elwin/heatmap/api/app"
	"github.com/stretchr/testify/suite"
)

var (
	update = flag.Bool("update", false, "update the golden files of this test")
	port   = ":3000"
	host   = "http://localhost" + port
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

		// xx.Strava.RememberID = "***REMOVED***"
		// xx.Strava.RememberToken = "***REMOVED***"

		token, err := FetchToken(ctx, conf, xx.Strava.RememberID, xx.Strava.RememberToken)
		suite.Require().NoError(err)

		suite.client = NewClient(ctx, conf, token)
	}
}

func (suite *StravaTestSuite) Test_activities() {
	ctx := context.Background()

	activities, err := suite.client.activities().all(ctx)
	suite.Require().NoError(err)

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
