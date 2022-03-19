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
	"os"

	"github.com/gmemstr/pogo/router"
)

// Main function that defines routes
func main() {

	// Check if this is the first time Pogo has been run
	// with a lockfile
	if _, err := os.Stat("assets/.lock"); err != nil {
		fmt.Println("This looks like your first time running Pogo, give me a second to set myself up.")
		Setup()
	}

	// Start the watch() function in generate_rss.go, which
	// watches for file changes and regenerates the feed
	go watch()

	// Define routes
	// We're live
	r := router.Init()
	fmt.Println("Your Pogo instance is live on port :3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
