package main

import (
	"io/ioutil"
	"net/http"
	// "fmt"
	"encoding/json"
	"strings"
)

type NewConfig struct {
	Name        string
	Host        string
	Email       string
	Description string
	Image       string
	PodcastUrl  string
}

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
		r.ParseMultipartForm(32 << 20)

		// Parse form and convert to JSON
		cnf := NewConfig{
			strings.Join(r.Form["podcastname"], ""),        // Podcast name
			strings.Join(r.Form["podcasthost"], ""),        // Podcast host
			strings.Join(r.Form["podcastemail"], ""),       // Podcast host email
			strings.Join(r.Form["podcastdescription"], ""), // Podcast Description
			"", // Podcast image
			"", // Podcast location
		}

		b, err := json.Marshal(cnf)
		if err != nil {
			panic(err)
		}

		ioutil.WriteFile("config.json", b, 0644)
		w.Write([]byte("Done"))
	}
}
