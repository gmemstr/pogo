package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gmemstr/feeds"
	"github.com/spf13/viper"
)

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
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	now := time.Now()
	files, err := ioutil.ReadDir("podcasts")
	if err != nil {
		log.Fatal(err)
	}

	feed := &feeds.Feed{
		Title:       viper.GetString("Name"),
		Link:        &feeds.Link{Href: viper.GetString("PodcastUrl")},
		Description: "discussion about open source projects",
		Author:      &feeds.Author{Name: viper.GetString("Host"), Email:viper.GetString("Email")},
		Created:     now,
		Image:       &feeds.Image{Url: viper.GetString("Image")},
	}

	for _, file := range files {
		if strings.Contains(file.Name(), ".mp3") {
			fmt.Println(file.Name())
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
			feed.Add (
				&feeds.Item{
					Title:       title,
					Link:        &feeds.Link{Href: viper.GetString("PodcastUrl") + "/download/" + file.Name(), Length: "100", Type: "audio/mpeg"},
					Enclosure:   &feeds.Enclosure{Url: viper.GetString("PodcastUrl") + "/download/" + file.Name(), Length: "100", Type: "audio/mpeg"},
					Description: string(description),
					Author:      &feeds.Author{Name: viper.GetString("Host"), Email: viper.GetString("Email")},
					Created:     date,
				},
			)
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
	// fmt.Println(rss)
	rss_byte := []byte(rss)
	ioutil.WriteFile("feed.rss", rss_byte, 0644)
	json_byte := []byte(json)
	ioutil.WriteFile("feed.json", json_byte, 0644)
}
