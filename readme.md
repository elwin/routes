# Routes
This CLI application currently does two things:
1. Download all activities from Strava
2. Generate a poster using those activities, as seen below

## Installation
Requires a working go installation, including correctly set up `$GOPATH` variable:
```shell
go install github.com/elwin/routes/cmd/routes
```

## Usage
Appropriate client ID and secret can be created in your Strava account.

```shell
routes strava --client_id 123 --client_secret xyz --output activities.json
routes poster --input activities.json
```



![Mockup](resources/mockup.jpg)
Original Photo by [Martin Péchy](https://unsplash.com/photos/iXHdGk8JVYU)