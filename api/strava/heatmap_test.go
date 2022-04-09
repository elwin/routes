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
)

type StravaTestSuite struct {
	suite.Suite
	conf   Config
	client *Client
}

func (suite *StravaTestSuite) SetupTest() {
	if suite.client == nil {
		ctx := context.Background()

		appConf, err := app.ReadConfig("../config.yml")
		suite.Require().NoError(err)
		conf := Config{
			ClientID:     appConf.Api.Strava.ClientID,
			ClientSecret: appConf.Api.Strava.ClientSecret,
			RedirectHost: appConf.Api.Host,
		}

		token, err := FetchToken(ctx, conf)
		suite.Require().NoError(err)

		suite.client = NewClient(ctx, conf, token)
	}
}

func (suite *StravaTestSuite) Test_something() {
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
		if len(activities) == 10 {
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
