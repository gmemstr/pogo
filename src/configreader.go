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
}

type User struct {
	Username string
	Hash string
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

func GetUser(username string) (usr string, pwd string) {
	d, err := ioutil.ReadFile("users.json")
	if err != nil {
		panic(err)
	}

	var u interface{}
	err = json.Unmarshal(d, &u)

	users := u.(map[string]interface{})
	for k, v := range users {
    	if k == username {
    		usr = k
    		pwd = v.(string)
    	}
	}
	return
}