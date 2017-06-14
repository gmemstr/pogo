package main

import (
    "fmt"
    "log"
    "io/ioutil"
    "net/http"
)
func rootHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/rss+xml")
    w.WriteHeader(http.StatusOK)
    data, err := ioutil.ReadFile("feed.rss")
    if err != nil {
        panic(err)
    }
    w.Header().Set("Content-Length", fmt.Sprint(len(data)))
    fmt.Fprint(w, string(data))
}

func main() {
    go watch()
    http.HandleFunc("/", rootHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
    fmt.Println("watching folder")
}