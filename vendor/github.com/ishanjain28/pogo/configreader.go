package main

import (
	"encoding/json"
	"io/ioutil"
)

// Configuration structure
type Config struct {
	Name        string
	Host        string
	Email       string
	Description string
	Image       string
	PodcastUrl  string
}

// Single use structure
type User struct {
	Username string
	Hash     string
}

// Read config file and make values accesible
func ReadConfig() (c Config, err error) {
	// Read config.json
	d, err := ioutil.ReadFile("assets/config/config.json")
	if err != nil {
		return
	}

	err = json.Unmarshal(d, &c)
	if err != nil {
		return
	}

	return
}

// Return single users username & passsword *hash*
func AuthUser(username string, password string) (authd bool) {
	// Read users json file
	d, err := ioutil.ReadFile("assets/config/users.json")
	if err != nil {
		panic(err)
	}

	var u interface{}
	err = json.Unmarshal(d, &u) // Unmarshal into interface

	// Iterate through map until we find matching username
	users := u.(map[string]interface{})
	for k, v := range users {
		if k == username && v.(string) == password {
			authd = true
		}
	}
	return // Returns whether user is authenticated or not
}
