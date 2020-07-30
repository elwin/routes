package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type config struct {
	Host   string `yaml:"host"`
	Strava struct {
		ID     string `yaml:"id"`
		Secret string `yaml:"secret"`
	} `yaml:"strava"`
}

func readConfig(path string) (config, error) {
	f, err := os.Open(path)
	if err != nil {
	    return config{}, err
	}

	var c config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&c); err != nil {
		return config{}, err
	}

	return c, nil
}