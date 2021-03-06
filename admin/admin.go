/* admin.go
 *
 * Here is where all the neccesary functions for managing episodes
 * live, e.g adding removing etc.
 */

package admin

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"

	"github.com/gmemstr/pogo/common"
)

type User struct {
	Id   int    `json:"id"`
	Dbun string `json:"username"`
	Dbrn string `json:"realname"`
	Dbem string `json:"email"`
}
type UserList struct {
	Users []User
}
type Config struct {
	Name        string
	Host        string
	Email       string
	Description string
	Image       string
	PodcastUrl  string
}

/*
 * The following is a set of admin commands
 * that the average user probably shouldn't be
 * able to have access to, mostly user management.
 */

func ConfigurationManager() common.Handler {

	return func(rc *common.RouterContext, w http.ResponseWriter, r *http.Request) *common.HTTPError {

		if r.Method == "GET" {
			return common.ReadAndServeFile("assets/config/config.json", w)
		}
		fmt.Println(r.Host)
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			return &common.HTTPError{
				Message:    err.Error(),
				StatusCode: http.StatusBadRequest,
			}
		}
		file, handler, err := r.FormFile("image")
		if err == nil {
			defer file.Close()
			fmt.Println(handler.Header)
			f, err := os.OpenFile("./assets/web/static/podcastimage.png", os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				w.Write([]byte("<script>window.location = '/admin#failed';</script>"))

				fmt.Println(err)
				return nil
			}
			defer f.Close()
			io.Copy(f, file)
		} else {
			// Do nothing, assume no image was uploaded :(
		}
		newconfig := Config{
			strings.Join(r.Form["feedname"], ""),
			strings.Join(r.Form["host"], ""),
			strings.Join(r.Form["email"], ""),
			strings.Join(r.Form["description"], ""),
			r.Host + "/assets/podcastimage.png",
			"http://" + r.Host,
		}
		newconfigjson, err := json.Marshal(newconfig)
		if err != nil {
			return &common.HTTPError{
				Message:    err.Error(),
				StatusCode: http.StatusBadRequest,
			}
		}
		err = ioutil.WriteFile("./assets/config/config.json", newconfigjson, 0644)
		if err != nil {
			return &common.HTTPError{
				Message:    err.Error(),
				StatusCode: http.StatusBadRequest,
			}
		}
		w.Write([]byte("<script>window.location = '/admin#/Settings#saved';</script>"))
		return nil
	}

}

func AddUser() common.Handler {

	return func(rc *common.RouterContext, w http.ResponseWriter, r *http.Request) *common.HTTPError {

		db, err := sql.Open("sqlite3", "assets/config/users.db")
		if err != nil {
			return &common.HTTPError{
				Message:    fmt.Sprintf("error opening sqlite3 file: %v", err),
				StatusCode: http.StatusInternalServerError,
			}
		}
		statement, err := db.Prepare("INSERT INTO users(username,hash,realname,email,permissions) VALUES (?,?,?,?,?)")
		if err != nil {
			return &common.HTTPError{
				Message:    fmt.Sprintf("error preparing sqlite3 statement: %v", err),
				StatusCode: http.StatusInternalServerError,
			}
		}

		err = r.ParseMultipartForm(32 << 20)
		if err != nil {
			return &common.HTTPError{
				Message:    err.Error(),
				StatusCode: http.StatusBadRequest,
			}
		}

		username := strings.Join(r.Form["username"], "")
		password := strings.Join(r.Form["password"], "")
		realname := strings.Join(r.Form["realname"], "")
		email := strings.Join(r.Form["email"], "")
		permissions := strings.Join(r.Form["permissions"], "")

		hash, err := bcrypt.GenerateFromPassword([]byte(password), 4)

		_, err = statement.Exec(username, hash, realname, email, permissions)
		if err != nil {
			return &common.HTTPError{
				Message:    fmt.Sprintf("error executing sqlite3 statement: %v", err),
				StatusCode: http.StatusInternalServerError,
			}
		}
		w.Write([]byte("<script>window.location = '/admin#/users/added';</script>"))
		db.Close()
		return nil
	}

}

func EditUser() common.Handler {

	return func(rc *common.RouterContext, w http.ResponseWriter, r *http.Request) *common.HTTPError {
		db, err := sql.Open("sqlite3", "assets/config/users.db")

		if err != nil {
			return &common.HTTPError{
				Message:    fmt.Sprintf("error in reading user database: %v", err),
				StatusCode: http.StatusInternalServerError,
			}
		}

		err = r.ParseMultipartForm(32 << 20)
		if err != nil {
			return &common.HTTPError{
				Message:    err.Error(),
				StatusCode: http.StatusBadRequest,
			}
		}
		id := strings.Join(r.Form["id"], "")
		username := strings.Join(r.Form["username"], "")
		password := strings.Join(r.Form["oldpw"], "")
		newpassword := strings.Join(r.Form["newpw1"], "")
		realname := strings.Join(r.Form["realname"], "")
		email := strings.Join(r.Form["email"], "")
		permissions := strings.Join(r.Form["permissions"], "")
		pwhash, err := bcrypt.GenerateFromPassword([]byte(password), 4)

		statement, err := db.Prepare("UPDATE users SET username=?, hash=?, realname=?, email=?, permissions=? WHERE id=?")
		if err != nil {
			return &common.HTTPError{
				Message:    fmt.Sprintf("error preparing sqlite3 statement: %v", err),
				StatusCode: http.StatusInternalServerError,
			}
		}

		pwstatement, err := db.Prepare("SELECT hash FROM users WHERE id=?")
		if err != nil {
			return &common.HTTPError{
				Message:    fmt.Sprintf("error preparing sqlite3 statement: %v", err),
				StatusCode: http.StatusInternalServerError,
			}
		}

		tmp, err := pwstatement.Query(id)
		if err != nil {
			return &common.HTTPError{
				Message:    fmt.Sprintf("error executing sqlite3 statement: %v", err),
				StatusCode: http.StatusInternalServerError,
			}
		}

		var hash []byte

		for tmp.Next() {
			err = tmp.Scan(&hash)
			if err != nil {
				return &common.HTTPError{
					Message:    fmt.Sprintf("error executing sqlite3 statement: %v", err),
					StatusCode: http.StatusInternalServerError,
				}
			}
		}
		fmt.Println(hash)
		if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
			fmt.Println("Passwords do not match")
			w.Write([]byte("<script>window.location = '/admin#/users/editerror';</script>"))
			db.Close()

			return nil
		}

		if newpassword != "" {
			pwhash, err = bcrypt.GenerateFromPassword([]byte(newpassword), 4)
		}

		_, err = statement.Exec(username, pwhash, realname, email, id, permissions)
		if err != nil {
			return &common.HTTPError{
				Message:    fmt.Sprintf("error executing sqlite3 statement: %v", err),
				StatusCode: http.StatusInternalServerError,
			}
		}
		w.Write([]byte("<script>window.location = '/admin#/users/edited';</script>"))
		db.Close()

		return nil
	}
}

func DeleteUser() common.Handler {

	return func(rc *common.RouterContext, w http.ResponseWriter, r *http.Request) *common.HTTPError {

		db, err := sql.Open("sqlite3", "assets/config/users.db")
		if err != nil {
			return &common.HTTPError{
				Message:    fmt.Sprintf("error opening sqlite3 file: %v", err),
				StatusCode: http.StatusInternalServerError,
			}
		}
		statement, err := db.Prepare("DELETE FROM users WHERE id=?")
		if err != nil {
			return &common.HTTPError{
				Message:    fmt.Sprintf("error preparing sqlite3 statement: %v", err),
				StatusCode: http.StatusInternalServerError,
			}
		}

		if err != nil {
			return &common.HTTPError{
				Message:    err.Error(),
				StatusCode: http.StatusBadRequest,
			}
		}
		vars := mux.Vars(r)
		id := vars["id"]
		if id == "1" {
			w.Write([]byte("<script>window.location = '/admin#/msg/Cannot%20Delete%20Admin%20User';</script>"))
			db.Close()
			return nil
		}

		_, err = statement.Exec(id)
		if err != nil {
			return &common.HTTPError{
				Message:    fmt.Sprintf("error executing sqlite3 statement: %v", err),
				StatusCode: http.StatusInternalServerError,
			}
		}
		w.Write([]byte("<script>window.location = '/admin#/msg/Deleted%20User';</script>"))
		db.Close()
		return nil
	}
}

func ListUsers() common.Handler {

	return func(rc *common.RouterContext, w http.ResponseWriter, r *http.Request) *common.HTTPError {

		db, err := sql.Open("sqlite3", "assets/config/users.db")

		if err != nil {
			return &common.HTTPError{
				Message:    fmt.Sprintf("error in reading user database: %v", err),
				StatusCode: http.StatusInternalServerError,
			}
		}
		// NEVER SELECT hash ENTRY
		statement, err := db.Prepare("SELECT id,username,realname,email FROM users")
		if err != nil {
			return &common.HTTPError{
				Message:    fmt.Sprintf("error in reading user database: %v", err),
				StatusCode: http.StatusInternalServerError,
			}
		}

		rows, err := statement.Query()
		if err != nil {
			return &common.HTTPError{
				Message:    fmt.Sprintf("error in executing user SELECT: %v", err),
				StatusCode: http.StatusInternalServerError,
			}
		}
		res := []User{}

		for rows.Next() {
			var u User
			err := rows.Scan(&u.Id, &u.Dbun, &u.Dbrn, &u.Dbem)
			if err != nil {
				return &common.HTTPError{
					Message:    fmt.Sprintf("error in decoding sql data", err),
					StatusCode: http.StatusBadRequest,
				}
			}
			res = append(res, u)
		}
		fin, err := json.Marshal(res)
		w.Write(fin)
		db.Close()

		return nil
	}
}

/*************************************
 * End of "sensitive" admin functions
 * ***********************************/

// Write custom CSS to disk or send it back to the client if GET

func CustomCss() common.Handler {

	return func(rc *common.RouterContext, w http.ResponseWriter, r *http.Request) *common.HTTPError {

		if r.Method == "GET" {
			return common.ReadAndServeFile("assets/web/static/custom.css", w)
		}

		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			return &common.HTTPError{
				Message:    err.Error(),
				StatusCode: http.StatusBadRequest,
			}
		}

		css := strings.Join(r.Form["css"], "")

		filename := "custom.css"

		err = ioutil.WriteFile("./assets/web/static/"+filename, []byte(css), 0644)

		if err != nil {
			w.Write([]byte("<script>window.location = '/admin#failed';</script>"))
			panic(err)
		} else {
			w.Write([]byte("<script>window.location = '/admin#cssupdated';</script>"))
		}
		return nil
	}
}

func EditEpisode() common.Handler {
	return func(rc *common.RouterContext, w http.ResponseWriter, r *http.Request) *common.HTTPError {
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			return &common.HTTPError{
				Message:    err.Error(),
				StatusCode: http.StatusBadRequest,
			}
		}

		PreviousFilename := strings.Join(r.Form["previousfilename"], "")

		date := strings.Join(r.Form["date"], "")
		title := strings.Join(r.Form["title"], "")

		name := fmt.Sprintf("%v_%v", date, title)
		filename := "./podcasts/" + name + ".mp3"
		shownotes := "./podcasts/" + name + "_SHOWNOTES.md"
		fmt.Println(filename)
		description := strings.Join(r.Form["description"], "")

		if "./podcasts"+PreviousFilename+".mp3" != filename {
			err = os.Rename("./podcasts/"+PreviousFilename+".mp3", filename)
			if err != nil {
				return &common.HTTPError{
					Message:    err.Error(),
					StatusCode: http.StatusBadRequest,
				}
			}
			err = os.Rename("./podcasts/"+PreviousFilename+"_SHOWNOTES.md", shownotes)
			if err != nil {
				return &common.HTTPError{
					Message:    err.Error(),
					StatusCode: http.StatusBadRequest,
				}
			}
		}
		err = ioutil.WriteFile(shownotes, []byte(description), 0644)
		if err != nil {
			return &common.HTTPError{
				Message:    err.Error(),
				StatusCode: http.StatusBadRequest,
			}
		}
		w.Write([]byte("<script>window.location = '/admin#/msg/Episode%20Published!';</script>"))
		return nil
	}
}

func CreateEpisode() common.Handler {
	return func(rc *common.RouterContext, w http.ResponseWriter, r *http.Request) *common.HTTPError {
		d, err := ioutil.ReadFile("assets/config/config.json")
		if err != nil {
			panic(err)
		}
		var config Config
		err = json.Unmarshal(d, &config)
		if err != nil {
			panic(err)
		}
		err = r.ParseMultipartForm(32 << 20)
		if err != nil {
			return &common.HTTPError{
				Message:    err.Error(),
				StatusCode: http.StatusBadRequest,
			}
		}

		// Build filename for episode
		date := strings.Join(r.Form["date"], "")
		title := strings.Join(r.Form["title"], "")
		db, err := sql.Open("sqlite3", "assets/config/users.db")
		defer db.Close()
		author := config.Host

		if err != nil {
			fmt.Println("Error getting user from database", err)
		}

		statement, err := db.Prepare("SELECT realname FROM users WHERE username=?")
		if err != nil {
			fmt.Println("Error getting user from database", err)
		}

		rows, err := statement.Query(rc.User.Username)
		if err != nil {
			fmt.Println("Error getting user from database", err)
		}

		var realname string
		for rows.Next() {
			err = rows.Scan(&realname)
			if err != nil {
				fmt.Println("Error getting user from database", err)
			}
			author = realname
		}
		name := fmt.Sprintf("%v_%v", date, title)
		filename := name + ".mp3"
		shownotes := name + "_SHOWNOTES.md"
		description := author + "\n" + strings.Join(r.Form["description"], "")
		// Finish building filenames

		err = ioutil.WriteFile("./podcasts/"+shownotes, []byte(description), 0644)
		if err != nil {
			w.Write([]byte("<script>window.location = '/admin#failed';</script>"))
			fmt.Println(err)
		}

		file, handler, err := r.FormFile("file")
		if err != nil {
			w.Write([]byte("<script>window.location = '/admin#failed';</script>"))

			fmt.Println(err)
			return nil
		}
		defer file.Close()
		fmt.Println(handler.Header)
		f, err := os.OpenFile("./podcasts/"+filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			w.Write([]byte("<script>window.location = '/admin#failed';</script>"))

			fmt.Println(err)
			return nil
		}
		defer f.Close()
		io.Copy(f, file)
		w.Write([]byte("<script>window.location = '/admin#published';</script>"))

		return nil
	}
}

func RemoveEpisode() common.Handler {
	return func(rc *common.RouterContext, w http.ResponseWriter, r *http.Request) *common.HTTPError {
		// Episode should be the full MP3 filename
		// Remove MP3 first
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			return &common.HTTPError{
				Message:    err.Error(),
				StatusCode: http.StatusBadRequest,
			}
		}
		episode := strings.Join(r.Form["episode"], "")
		os.Remove(episode)
		sn := strings.Replace(episode, ".mp3", "_SHOWNOTES.md", 2)
		os.Remove(sn)

		return nil
	}
}
