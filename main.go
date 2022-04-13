package main

import (
	"github.com/joho/godotenv"
	"wildlife/internal/log"
	"wildlife/internal/server"
	"wildlife/internal/server/controller"
	"wildlife/internal/server/tensor"
)

const VERSION = "0.0.1"
const OnnxModel = "assets/poisonOak.onnx"

func main() {
	log.Logf("Started v%s", VERSION)

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Errf("Error loading .env file: %s", err)
	}

	// Initialize the CNN model
	err = tensor.BuildModel(OnnxModel, false)
	if err != nil {
		log.Errf("Error initializing model: %s", err)
		return
	}

	// Initialize the controllers with the database
	err = controller.InitController()
	if err != nil {
		log.Errf("Error initializing controller: %s", err)
		return
	}

	// Start the web server
	err = server.Start()
	if err != nil {
		log.Errf("Error starting server: %s", err)
		return
	}
	log.Logf("Exiting.")
}
