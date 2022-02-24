package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"wildlife/internal/log"
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
	srv := http.Server{
		Addr:    ":3060",
		Handler: router,
	}
	// Start hosting the server
	if err := srv.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
