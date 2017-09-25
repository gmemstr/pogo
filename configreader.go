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
func ReadConfig() Config {
	// Read config.json
	d, err := ioutil.ReadFile("assets/config/config.json")
	if err != nil {
		panic(err)
	}

	var c Config // Unmarshal json
	err = json.Unmarshal(d, &c)
	if err != nil {
		panic(err)
	}

	return c
}

// Return single users username & passsword *hash*
func GetUser(username string) (usr string, pwd string) {
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
		if k == username {
			usr = k
			pwd = v.(string)
		}
	}
	return // Returns k & v values, aka username and password *hash*
}
