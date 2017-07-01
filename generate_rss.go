package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gmemstr/feeds"
	"github.com/Tkanos/gonfig"
)

type Configuration struct {
	Name	   string
	Host	   string
	Email	   string
	Image      string
	PodcastUrl string
}

func watch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
					generate_rss()
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("podcasts/")
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func generate_rss() {
	configuration := Configuration{}
	err := gonfig.GetConf("config.json", &configuration)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(configuration)
	now := time.Now()
	files, err := ioutil.ReadDir("podcasts")
	if err != nil {
		log.Fatal(err)
	}

	feed := &feeds.Feed{
		Title:       configuration.Name,
		Link:        &feeds.Link{Href: configuration.PodcastUrl},
		Description: "discussion about open source projects",
		Author:      &feeds.Author{Name: configuration.Host, Email: configuration.Email},
		Created:     now,
		Image:       &feeds.Image{Url: configuration.Image},
	}

	for _, file := range files {
		if strings.Contains(file.Name(), ".mp3") {
			s := strings.Split(file.Name(), "_")
			t := strings.Split(s[1], ".")
			title := t[0]
			description,err := ioutil.ReadFile("podcasts/" + strings.Replace(file.Name(), ".mp3", "_SHOWNOTES.md", 2))
			if err != nil {
		        log.Fatal(err)
		    }
			date, err := time.Parse("2006-01-02", s[0])
			if err != nil {
				log.Fatal(err)
			}
			feed.Items = []*feeds.Item{
				&feeds.Item{
					Title:       title,
					Link:        &feeds.Link{Href: configuration.PodcastUrl + "/download/" + file.Name(), Length: "100", Type: "audio/mpeg"},
					Enclosure:   &feeds.Enclosure{Url: configuration.PodcastUrl + "/download/" + file.Name(), Length: "100", Type: "audio/mpeg"},
					Description: string(description),
					Author:      &feeds.Author{Name: configuration.Host, Email: configuration.Email},
					Created:     date,
				},
			}
		}
	}
	rss, err := feed.ToRss()
	if err != nil {
		log.Fatal(err)
	}
	json, err := feed.ToJSON()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(rss)
	rss_byte := []byte(rss)
	ioutil.WriteFile("feed.rss", rss_byte, 0644)
	json_byte := []byte(json)
	ioutil.WriteFile("feed.json", json_byte, 0644)
}
