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
	var c Config

	if err := ReadYamlConfig(path, &c); err != nil {
		return Config{}, err
	}

	return c, nil
}

func ReadYamlConfig(path string, x interface{}) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return yaml.NewDecoder(f).Decode(x)
}