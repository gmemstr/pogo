package main

import (
	"io/ioutil"
	"net/http"
)

// Serve setup.html and config parameters
func ServeSetup(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data, err := ioutil.ReadFile("assets/setup.html")

		if err != nil {
			panic(err)
		}
		w.Write(data)
	}
	if r.Method == "POST" {
		
	}
}

// Write JSON config to file
func WriteConfig() {

}

// Connect to SQL and create admin user
func CreateAdmin() {

}
