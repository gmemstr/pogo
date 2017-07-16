package main

import (
	"net/http"
	"fmt"
	"strings"
	"io/ioutil"
	"io"
	"os"
)

func CreateEpisode(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20)
		date := strings.Join(r.Form["date"], "")
		title := strings.Join(r.Form["title"], "")

		name :=  fmt.Sprintf("%v_%v", date, title)
		filename := name + ".mp3"
		shownotes := name + "_SHOWNOTES.md"
		fmt.Println(name)
		description := strings.Join(r.Form["description"], "")
		fmt.Println(description)

		err := ioutil.WriteFile("./podcasts/" + shownotes, []byte(description), 0644)
	    if err != nil {
	        panic(err)
	    }

		file, handler, err := r.FormFile("file")
	    if err != nil {
	        fmt.Println(err)
	        return
	    }
	    defer file.Close()
	    fmt.Fprintf(w, "%v", handler.Header)
	    f, err := os.OpenFile("./podcasts/"+filename, os.O_WRONLY|os.O_CREATE, 0666)
	    if err != nil {
	        fmt.Println(err)
	        return
	    }
	    defer f.Close()
	    io.Copy(f, file)
	}
}