/* generate_rss.go
 *
 * This file contains functions for monitoring for file changes and
 * regenerating the RSS feed accordingly, pulling in shownote files
 * and configuration parameters
 */

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

// Watch folder for changes, called from webserver.go
func watch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)

	// Call func asynchronously
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				// log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					// log.Println("modified file:", event.Name)
					log.Println("File up(load/date)ed: ", event.Name)
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
	err = watcher.Add("config.json")
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

// Called when a file has been created / changed, uses gorilla feeds
// fork to add items to feed object
func generate_rss() {
	config := ReadConfig()
	now := time.Now()
	files, err := ioutil.ReadDir("podcasts")
	if err != nil {
		log.Fatal(err)
	}

	podcasturl := config.PodcastUrl
	feed := &feeds.Feed{
		Title:       config.Name,
		Link:        &feeds.Link{Href: podcasturl},
		Description: config.Description,
		Author:      &feeds.Author{Name: config.Host, Email: config.Email},
		Created:     now,
		Image:       &feeds.Image{Url: config.Image},
	}

	for _, file := range files {
		if strings.Contains(file.Name(), ".mp3") {
			s := strings.Split(file.Name(), "_")
			t := strings.Split(s[1], ".")
			title := t[0]
			description, err := ioutil.ReadFile("podcasts/" + strings.Replace(file.Name(), ".mp3", "_SHOWNOTES.md", 2))
			if err != nil {
				log.Fatal(err)
			}
			date, err := time.Parse("2006-01-02", s[0])
			if err != nil {
				log.Fatal(err)
			}
			size := fmt.Sprintf("%d", file.Size())
			link := podcasturl + "/download/" + file.Name()
			feed.Add(
				&feeds.Item{
					Title:       title,
					Link:        &feeds.Link{Href: link, Length: size, Type: "audio/mpeg"},
					Enclosure:   &feeds.Enclosure{Url: link, Length: size, Type: "audio/mpeg"},
					Description: string(description),
					Author:      &feeds.Author{Name: config.Host, Email: config.Email},
					Created:     date,
				},
			)
		}
	}

	// Translate the feed to both RSS and JSON,
	// RSS for readers and JSON for frontend (& API I guess)
	rss, err := feed.ToRss()
	if err != nil {
		log.Fatal(err)
	}
	json, err := feed.ToJSON()
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(rss)

	// Write to files as neccesary
	rss_byte := []byte(rss)
	ioutil.WriteFile("feed.rss", rss_byte, 0644)

	json_byte := []byte(json)
	ioutil.WriteFile("feed.json", json_byte, 0644)
}
