/* admin.go
 *
 * Here is where all the neccesary functions for managing episodes
 * live, e.g adding removing etc.
 */

package admin

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gmemstr/pogo/common"
)

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

		if ("./podcasts" + PreviousFilename + ".mp3" != filename) {
			err = os.Rename("./podcasts/" + PreviousFilename + ".mp3", filename)
			if err != nil {
				return &common.HTTPError{
					Message:    err.Error(),
					StatusCode: http.StatusBadRequest,
				}
			}
			err = os.Rename("./podcasts/" + PreviousFilename + "_SHOWNOTES.md", shownotes)
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
		w.Write([]byte("<script>window.location = '/admin#published';</script>"))
		return nil
	}
}

func CreateEpisode() common.Handler {
	return func(rc *common.RouterContext, w http.ResponseWriter, r *http.Request) *common.HTTPError {
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			return &common.HTTPError{
				Message:    err.Error(),
				StatusCode: http.StatusBadRequest,
			}
		}

		// Build filename for episode
		date := strings.Join(r.Form["date"], "")
		title := strings.Join(r.Form["title"], "")

		name := fmt.Sprintf("%v_%v", date, title)
		filename := name + ".mp3"
		shownotes := name + "_SHOWNOTES.md"
		description := strings.Join(r.Form["description"], "")
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
