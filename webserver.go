/* webserver.go
 *
 * This is the webserver handler for Pogo, and handles
 * all incoming connections, including authentication.
 */

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ishanjain28/pogo/router"
)

// Authenticate user using basic webserver authentication
// @TODO: Replace this with a different for of _more secure_
// authentication that we can POST to instead.
/*
 * Code from stackoverflow by user Timmmm
 * https://stackoverflow.com/questions/21936332/idiomatic-way-of-requiring-http-basic-auth-in-go/39591234#39591234
 */
// func BasicAuth(handler http.HandlerFunc) http.HandlerFunc {

// 	return func(w http.ResponseWriter, r *http.Request) {
// 		realm := "Login to Pogo admin interface"
// 		user, pass, ok := r.BasicAuth()

// 		if !AuthUser(user, pass) || !ok {
// 			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
// 			w.WriteHeader(401)
// 			w.Write([]byte("Unauthorised.\n"))
// 			return
// 		}
// 		handler(w, r)
// 	}
// }

// // Handler for serving up admin page
// func AdminHandler(w http.ResponseWriter, r *http.Request) {
// 	data, err := ioutil.ReadFile("assets/web/admin.html")

// 	if err == nil {
// 		w.Write(data)
// 	} else {
// 		w.WriteHeader(500)
// 		w.Write([]byte("500 Something went wrong - " + http.StatusText(500)))
// 	}
// }

// Main function that defines routes
func main() {
	// Start the watch() function in generate_rss.go, which
	// watches for file changes and regenerates the feed
	// go watch()

	// Define routes
	// We're live
	r := router.Init()
	fmt.Println("Listening on port :3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
