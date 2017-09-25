package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Name          string
	Host          string
	Email         string
	Description   string
	Image         string
	PodcastUrl    string
	AdminUsername string
	AdminPassword string
}

func ReadConfig() Config {
	d, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	var c Config
	err = json.Unmarshal(d, &c)
	if err != nil {
		panic(err)
	}

	return c
}
