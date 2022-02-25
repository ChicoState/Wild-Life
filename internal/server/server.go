package server

import (
	"net/http"
	"os"
	"wildlife/internal/log"
	"wildlife/internal/server/controller"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Start will initialize and start hosting the Rest API.
func Start() error {
	// Initialize a new router
	router := chi.NewRouter()
	// Middlewares
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	// Overwrite default logger with internal logger
	router.Use(log.Middleware)
	// Status middleware
	router.Use(middleware.Heartbeat("/status"))
	// Route requests to individual routers
	router.Route("/upload", uploadRouter)
	// Create a new http Server object

	// Test db add/remove
	if os.Getenv("TEST_USER_AR_WEB") == "add" || os.Getenv("TEST_USER_AR_WEB") == "all" {
		router.Get("/test/user/add", controller.TestDBAdd)
	}

	if os.Getenv("TEST_USER_AR_WEB") == "remove" || os.Getenv("TEST_USER_AR_WEB") == "all" {
		router.Get("/test/user/remove", controller.TestDBRemove)
	}

	srv := http.Server{
		Addr:    ":3060",
		Handler: router,
	}
	log.Logf("Starting server on port localhost:3060")
	// Start hosting the server
	if err := srv.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
