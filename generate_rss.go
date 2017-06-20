package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gmemstr/feeds"
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
	now := time.Now()
	files, err := ioutil.ReadDir("podcasts")
	if err != nil {
		log.Fatal(err)
	}

	feed := &feeds.Feed{
		Title:       "Git Galaxy Stargazers",
		Link:        &feeds.Link{Href: "https://gitgalaxy.com"},
		Description: "discussion about open source projects",
		Author:      &feeds.Author{Name: "Gabriel Simmer", Email: "gabriel@gitgalaxy.com"},
		Created:     now,
		Image:       &feeds.Image{Url: "https://podcast.gitgalaxy.com/assets/podcast_image.png"},
	}

	for _, file := range files {
		s := strings.Split(file.Name(), "_")
		t := strings.Split(s[1], ".")
		title := t[0]
		date, err := time.Parse("2006-01-02", s[0])
		if err != nil {
			log.Fatal(err)
		}
		feed.Items = []*feeds.Item{
			&feeds.Item{
				Title:       title,
				Link:        &feeds.Link{Href: "https://podcast.gitgalaxy.com/download/" + file.Name(), Length: "100", Type: "audio/mpeg"},
				Enclosure:   &feeds.Enclosure{Url: "https://podcast.gitgalaxy.com/download/" + file.Name(), Length: "100", Type: "audio/mpeg"},
				Description: "Hello, World!",
				Author:      &feeds.Author{Name: "Gabriel Simmer", Email: "gabriel@gitgalaxy.com"},
				Created:     date,
			},
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
