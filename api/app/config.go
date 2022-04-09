package app

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Api struct {
		Host             string `yaml:"host"`
		SessionDirectory string `yaml:"sessions"`
		Debug            bool   `yaml:"debug"`
		Strava           struct {
			ClientID     string `yaml:"id"`
			ClientSecret string `yaml:"secret"`
		} `yaml:"strava"`
	} `yaml:"api"`
}

func ReadConfig(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}

	var c Config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&c); err != nil {
		return Config{}, err
	}

	return c, nil
}
