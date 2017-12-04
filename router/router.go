package router

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gmemstr/pogo/admin"
	"github.com/gmemstr/pogo/auth"
	"github.com/gmemstr/pogo/common"
	"github.com/gorilla/mux"
)

type NewConfig struct {
	Name        string
	Host        string
	Email       string
	Description string
	Image       string
	PodcastURL  string
}

// Handle takes multiple Handler and executes them in a serial order starting from first to last.
// In case, Any middle ware returns an error, The error is logged to console and sent to the user, Middlewares further up in chain are not executed.
func Handle(handlers ...common.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		rc := &common.RouterContext{}
		for _, handler := range handlers {
			err := handler(rc, w, r)
			if err != nil {
				log.Printf("%v", err)

				w.Write([]byte(http.StatusText(err.StatusCode)))

				return
			}
		}
	})
}

func Init() *mux.Router {

	r := mux.NewRouter()

	// "Static" paths
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/web/static"))))
	r.PathPrefix("/download/").Handler(http.StripPrefix("/download/", http.FileServer(http.Dir("podcasts"))))

	// Paths that require specific handlers
	r.Handle("/", Handle(
		rootHandler(),
	)).Methods("GET")

	r.Handle("/rss", Handle(
		rootHandler(),
	)).Methods("GET")

	r.Handle("/json", Handle(
		rootHandler(),
	)).Methods("GET")

	// RequireAuthorization() handles authentication
	// and takes a single argument for permission level.
	// 0 any user, 1 most users, 2 only admin users
	r.Handle("/admin", Handle(
		auth.RequireAuthorization(0),
		adminHandler(),
	)).Methods("GET", "POST")

	r.Handle("/login", Handle(
		loginHandler(),
	)).Methods("GET", "POST")

	r.Handle("/admin/publish", Handle(
		auth.RequireAuthorization(0),
		admin.CreateEpisode(),
	)).Methods("POST")

	r.Handle("/admin/edituser", Handle(
		auth.RequireAuthorization(2),
		admin.EditUser(),
	)).Methods("POST")

	r.Handle("/admin/newuser", Handle(
		auth.RequireAuthorization(2),
		admin.AddUser(),
	)).Methods("POST")
	r.Handle("/admin/deleteuser/{id}", Handle(
		auth.RequireAuthorization(2),
		admin.DeleteUser(),
	)).Methods("GET")
	r.Handle("/admin/edit", Handle(
		auth.RequireAuthorization(1),
		admin.EditEpisode(),
	)).Methods("POST")

	r.Handle("/admin/delete", Handle(
		auth.RequireAuthorization(1),
		admin.RemoveEpisode(),
	)).Methods("GET")

	r.Handle("/admin/css", Handle(
		auth.RequireAuthorization(1),
		admin.CustomCss(),
	)).Methods("GET", "POST")

	r.Handle("/admin/adduser", Handle(
		auth.RequireAuthorization(2),
		admin.AddUser(),
	)).Methods("POST")

	r.Handle("/admin/listusers", Handle(
		auth.RequireAuthorization(1),
		admin.ListUsers(),
	)).Methods("GET")

	r.Handle("/setup", Handle(
		serveSetup(),
	)).Methods("GET", "POST")

	return r
}

func loginHandler() common.Handler {
	return func(rc *common.RouterContext, w http.ResponseWriter, r *http.Request) *common.HTTPError {
		db, err := sql.Open("sqlite3", "assets/config/users.db")

		if err != nil {
			return &common.HTTPError{
				Message:    fmt.Sprintf("error in reading user database: %v", err),
				StatusCode: http.StatusInternalServerError,
			}
		}

		statement, err := db.Prepare("SELECT * FROM users WHERE username=?")

		if _, err := auth.DecryptCookie(r); err == nil {
			http.Redirect(w, r, "/admin", http.StatusTemporaryRedirect)
			return nil
		}

		if r.Method == "GET" {
			w.Header().Set("Content-Type", "text/html")
			return common.ReadAndServeFile("assets/web/login.html", w)
		}

		err = r.ParseForm()
		if err != nil {
			return &common.HTTPError{
				Message:    fmt.Sprintf("error in parsing form: %v", err),
				StatusCode: http.StatusBadRequest,
			}
		}

		username := r.Form.Get("username")
		password := r.Form.Get("password")
		rows, err := statement.Query(username)

		if username == "" || password == "" {
			return &common.HTTPError{
				Message:    "username or password is empty",
				StatusCode: http.StatusBadRequest,
			}
		}
		var id int
		var dbun string
		var dbhsh string
		var dbrn string
		var dbem string
		var dbperm int
		for rows.Next() {
			err := rows.Scan(&id, &dbun, &dbhsh, &dbrn, &dbem, &dbperm)
			if err != nil {
				return &common.HTTPError{
					Message:    fmt.Sprintf("error in decoding sql data", err),
					StatusCode: http.StatusBadRequest,
				}
			}

		}
		// Create a cookie here because the credentials are correct
		if dbun == username && bcrypt.CompareHashAndPassword([]byte(dbhsh), []byte(password)) == nil {
			c, err := auth.CreateSession(&common.User{
				Username: username,
			})
			if err != nil {
				return &common.HTTPError{
					Message:    err.Error(),
					StatusCode: http.StatusInternalServerError,
				}
			}

			// r.AddCookie(c)
			w.Header().Add("Set-Cookie", c.String())
			// And now redirect the user to admin page
			http.Redirect(w, r, "/admin", http.StatusTemporaryRedirect)
			db.Close()
			return nil
		}

		return &common.HTTPError{
			Message:    "Invalid credentials!",
			StatusCode: http.StatusUnauthorized,
		}
	}
}

// Handles /, /feed and /json endpoints
func rootHandler() common.Handler {
	return func(rc *common.RouterContext, w http.ResponseWriter, r *http.Request) *common.HTTPError {

		var file string
		switch r.URL.Path {
		case "/rss":
			w.Header().Set("Content-Type", "application/rss+xml")
			file = "assets/web/feed.rss"
		case "/json":
			w.Header().Set("Content-Type", "application/json")
			file = "assets/web/feed.json"
		case "/":
			w.Header().Set("Content-Type", "text/html")
			file = "assets/web/index.html"
		default:
			return &common.HTTPError{
				Message:    fmt.Sprintf("%s: Not Found", r.URL.Path),
				StatusCode: http.StatusNotFound,
			}
		}

		return common.ReadAndServeFile(file, w)
	}
}

func adminHandler() common.Handler {
	return func(rc *common.RouterContext, w http.ResponseWriter, r *http.Request) *common.HTTPError {
		return common.ReadAndServeFile("assets/web/admin.html", w)
	}
}

// Serve setup.html and config parameters
func serveSetup() common.Handler {
	return func(rc *common.RouterContext, w http.ResponseWriter, r *http.Request) *common.HTTPError {
		if r.Method == "GET" {
			return common.ReadAndServeFile("assets/web/setup.html", w)
		}
		r.ParseMultipartForm(32 << 20)

		// Parse form and convert to JSON
		cnf := NewConfig{
			strings.Join(r.Form["podcastname"], ""),  // Podcast name
			strings.Join(r.Form["podcasthost"], ""),  // Podcast host
			strings.Join(r.Form["podcastemail"], ""), // Podcast host email
			"", // Podcast image
			"", // Podcast location
			"", // Podcast location
		}

		b, err := json.Marshal(cnf)
		if err != nil {
			panic(err)
		}

		ioutil.WriteFile("assets/config/config.json", b, 0644)
		w.Write([]byte("Done"))
		return nil
	}
}
