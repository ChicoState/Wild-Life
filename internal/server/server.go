package server

import (
	"net/http"
	"os"
	"wildlife/internal/log"

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

	router.Use(cors)
	// Overwrite default logger with internal logger
	router.Use(log.Middleware)
	// Status middleware
	router.Use(middleware.Heartbeat("/status"))
	// Route requests to individual routers
	router.Route("/upload", uploadRouter)
	router.Route("/sockets", socketRouter)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3060"
	}

	// Create a new http Server object
	srv := http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	log.Logf("Starting server on port localhost:%s", port)
	// Start hosting the server
	if err := srv.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

// cors provides allows cross-origin requests, only in development mode
func cors(next http.Handler) http.Handler {
	// Create http Handler function
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Serve request
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
