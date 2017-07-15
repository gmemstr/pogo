package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"crypto/subtle"

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

func JsonHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data, err := ioutil.ReadFile("feed.json")
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


/*
 * Code from stackoverflow by user Timmmm
 * https://stackoverflow.com/questions/21936332/idiomatic-way-of-requiring-http-basic-auth-in-go/39591234#39591234
*/
func BasicAuth(handler http.HandlerFunc, username, password, realm string) http.HandlerFunc {

    return func(w http.ResponseWriter, r *http.Request) {

        user, pass, ok := r.BasicAuth()

        if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
            w.Header().Set("WWW-Authenticate", `Basic realm="White Rabbit"`)
            w.WriteHeader(401)
            w.Write([]byte("Unauthorised.\n"))
            return
        }

        handler(w, r)
    }
}

func AdminHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("assets/admin.html")

	if err == nil {
		w.Write(data)
	} else {
		w.WriteHeader(404)
		w.Write([]byte("404 Something went wrong - " + http.StatusText(404)))
	}
}

func main() {
	go watch()
	r := mux.NewRouter()
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/static"))))
	r.PathPrefix("/download/").Handler(http.StripPrefix("/download/", http.FileServer(http.Dir("podcasts"))))
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/rss", RssHandler)
	r.HandleFunc("/json", JsonHandler)
	http.Handle("/", r)
	r.HandleFunc("/admin", BasicAuth(AdminHandler, "g", "password", "Login to White Rabbit admin interface"))
	log.Fatal(http.ListenAndServe(":8000", r))
}
