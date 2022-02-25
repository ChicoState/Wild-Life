package main

import (
	"wildlife/internal/log"
	"wildlife/internal/server"
)

const VERSION = "0.0.1"

func main() {
	log.Logf("Started v%s", VERSION)
	// Start the web server
	err := server.Start()
	if err != nil {
		return
	}
	log.Logf("Exiting.")
}
