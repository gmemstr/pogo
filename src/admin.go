/* admin.go
 * 
 * Here is where all the neccesary functions for managing episodes
 * live, e.g adding removing etc.
*/

package main

import (
	"net/http"
	"fmt"
	"strings"
	"io/ioutil"
	"io"
	"os"
)

func CustomCss(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20)
		css := strings.Join(r.Form["css"], "")

		filename := "custom.css"

		err := ioutil.WriteFile("./assets/static/" + filename, []byte(css), 0644)
		if err != nil { 
			w.Write([]byte("<script>window.location = '/admin#failed';</script>"))
			
			panic(err) 
		} else {
			w.Write([]byte("<script>window.location = '/admin#cssupdated';</script>")) 
		}
	} else {
		css,err := ioutil.ReadFile("./assets/static/custom.css")
		if err != nil {
			panic (err)
		} else {
			w.Write(css)
		}
	}
}

func CreateEpisode(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20)

		// Build filename for episode
		date := strings.Join(r.Form["date"], "")
		title := strings.Join(r.Form["title"], "")

		name :=  fmt.Sprintf("%v_%v", date, title)
		filename := name + ".mp3"
		shownotes := name + "_SHOWNOTES.md"
		fmt.Println(name)
		description := strings.Join(r.Form["description"], "")
		fmt.Println(description)
		// Finish building filenames 

		err := ioutil.WriteFile("./podcasts/" + shownotes, []byte(description), 0644)
	    if err != nil {
			w.Write([]byte("<script>window.location = '/admin#failed';</script>"))
	        panic(err)
	    }

		file, handler, err := r.FormFile("file")
	    if err != nil {
			w.Write([]byte("<script>window.location = '/admin#failed';</script>"))

	        fmt.Println(err)
	        return
	    }
	    defer file.Close()
	    fmt.Fprintf(w, "%v", handler.Header)
	    f, err := os.OpenFile("./podcasts/"+filename, os.O_WRONLY|os.O_CREATE, 0666)
	    if err != nil {
			w.Write([]byte("<script>window.location = '/admin#failed';</script>"))

	        fmt.Println(err)
	        return
	    }
	    defer f.Close()
	    io.Copy(f, file)
		w.Write([]byte("<script>window.location = '/admin#published';</script>"))     
	}
}

func RemoveEpisode(w http.ResponseWriter, r *http.Request) {
	// Episode should be the full MP3 filename
	// Remove MP3 first
	r.ParseMultipartForm(32 << 20)

	episode := strings.Join(r.Form["episode"],"")
	os.Remove(episode)
	sn := strings.Replace(episode, ".mp3", "_SHOWNOTES.md", 2)
	os.Remove(sn)
}