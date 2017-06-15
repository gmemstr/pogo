package main

import (
    "fmt"
    "log"
    "io/ioutil"
    "net/http"

    "github.com/gorilla/mux"
)
func RssHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/rss+xml")
    w.WriteHeader(http.StatusOK)
    data, err := ioutil.ReadFile("feed.rss")
    if err != nil {
        panic(err)
    }
    w.Header().Set("Content-Length", fmt.Sprint(len(data)))
    fmt.Fprint(w, string(data))
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
      data, err := ioutil.ReadFile("assets/index.html")

    if err == nil {
        w.Write(data)
    } else {
        w.WriteHeader(404)
        w.Write([]byte("404 Something went wrong - " + http.StatusText(404)))
    }
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Gorilla!\n"))
}

func main() {
    go watch()
    r := mux.NewRouter()
    r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/static"))))
    r.HandleFunc("/", HomeHandler)
    r.HandleFunc("/download/{episode}", DownloadHandler)
    r.HandleFunc("/rss", RssHandler)
    http.Handle("/", r)
    log.Fatal(http.ListenAndServe(":8000", r))
}