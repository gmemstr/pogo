/* generate_rss.go
 *
 * This file contains functions for monitoring for file changes and
 * regenerating the RSS feed accordingly, pulling in shownote files
 * and configuration parameters
 */

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/feeds"
)

type Config struct {
	Name        string
	Host        string
	Email       string
	Description string
	Image       string
	PodcastUrl  string
}

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
					GenerateRss()
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
	err = watcher.Add("assets/config/config.json")
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

// Iterate through podcasts directory and build feed
// object, then compile as json and rss and write to file
func GenerateRss() {
	d, err := ioutil.ReadFile("assets/config/config.json")
	if err != nil {
		panic(err)
	}
	var config Config
	err = json.Unmarshal(d, &config)
	if err != nil {
		panic(err)
	}

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
	i := 0
	for _, file := range files {
		if strings.Contains(file.Name(), ".mp3") {
			s := strings.Split(file.Name(), "_")
			t := strings.Split(s[1], ".")
			title := t[0]
			descfilelines := File2lines("podcasts/" + strings.Replace(file.Name(), ".mp3", "_SHOWNOTES.md", 2))
			author := descfilelines[0]
			description := descfilelines[1]
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
					Id:          strconv.Itoa(i),
					Title:       title,
					Link:        &feeds.Link{Href: link, Length: size, Type: "audio/mpeg"},
					Enclosure:   &feeds.Enclosure{Url: link, Length: size, Type: "audio/mpeg"},
					Description: description,
					Author:      &feeds.Author{Name: author, Email: config.Email}, // Replace with author in shownotes
					Created:     date,
				},
			)
			i = i + 1
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
	ioutil.WriteFile("assets/web/feed.rss", rss_byte, 0644)

	json_byte := []byte(json)
	ioutil.WriteFile("assets/web/feed.json", json_byte, 0644)
}

// From https://siongui.github.io/2016/04/06/go-readlines-from-file-or-string/
func File2lines(filePath string) []string {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	return lines
}
