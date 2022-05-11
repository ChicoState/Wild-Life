package server

import (
	"context"
	"net/http"
	"os"
	"wildlife/internal/log"
	"wildlife/internal/server/controller"
	"wildlife/internal/server/orchestrator"
	"wildlife/internal/test"

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
	orch, err := orchestrator.NewOrchestrator()
	if err != nil {
		return err
	}
	defer func(o *orchestrator.Orchestrator) {
		err = o.Close()
		if err != nil {
			log.Errf("Orchestrator deconstruction failed: %s", err)
		}
	}(orch)
	router.Use(OrchestratorContext(orch))
	router.Route("/upload", uploadRouter)
	router.Route("/sockets", socketRouter)

	// Test database
	if os.Getenv("TEST_USER_ARV") == "add" || os.Getenv("TEST_USER_ARV") == "remove" || os.Getenv("TEST_USER_ARV") == "all" {
		test.TestDB()
		// This test might mess with the already loaded users
	}

	// Test db add/remove
	if os.Getenv("TEST_USER_ARV_WEB") == "add" || os.Getenv("TEST_USER_ARV_WEB") == "all" {
		router.Get("/test/user/add", controller.TestDBAdd)
	}

	if os.Getenv("TEST_USER_ARV_WEB") == "remove" || os.Getenv("TEST_USER_ARV_WEB") == "all" {
		router.Get("/test/user/remove", controller.TestDBRemove)
	}
	if os.Getenv("TEST_USER_ARV_WEB") == "view" || os.Getenv("TEST_USER_ARV_WEB") == "all" {
		router.Get("/test/user/view", controller.TestDBView)
	}

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

func OrchestratorContext(orch *orchestrator.Orchestrator) func(n http.Handler) http.Handler {
	return func(n http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "orch", orch)
			n.ServeHTTP(w, r.WithContext(ctx))
		})
	}
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
