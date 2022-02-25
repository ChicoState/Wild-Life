package main

import (
	"os"
	"wildlife/internal/log"
	"wildlife/internal/server"
	"wildlife/internal/server/controller"
	"wildlife/internal/server/model"

	"github.com/joho/godotenv"
)

const VERSION = "0.0.1"

func main() {
	log.Logf("Started v%s", VERSION)
	// Start the web server
	godotenv.Load()

	if os.Getenv("DB_ACTIVE") == "true" {
		model.InitDB()
		controller.InitController()
	}

	err := server.Start()
	if err != nil {
		return
	}
	log.Logf("Exiting.")
}
