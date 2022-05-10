package main

import (
	"wildlife/internal/env"
	"wildlife/internal/log"
	"wildlife/internal/server"
	"wildlife/internal/server/controller"
	"wildlife/internal/server/orchestrator"
	"wildlife/internal/server/tensor"
)

//ONNX: Open Neural Network Exchange
//Allows us to train a model using python
//Save the model to a universal file
//Access it in go!
const (
	VERSION = "0.1.4"
	ASSETS  = "assets"
)

func main() {
	log.Logf("Started v%s", VERSION)

	// Load Environment
	err := env.Load()
	if err != nil {
		log.Errf("Environment initialization failed: %s", err)
	}

	// Initialize the CNN model
	err = tensor.BuildModel(ASSETS, false)
	if err != nil {
		log.Errf("Onnx Model initialization failed: %s", err)
		return
	}

	var o *orchestrator.Orchestrator
	o, err = orchestrator.NewOrchestrator()
	if err != nil {
		log.Errf("Orchestrator initialization failed: %s", err)
		return
	}
	defer func(o *orchestrator.Orchestrator) {
		err = o.Close()
		if err != nil {
			log.Errf("Orchestrator deconstruction failed: %s", err)
		}
	}(o)

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
