package controller

import (
	"os"
	"wildlife/internal/log"
	"wildlife/internal/server/model"
)

var (
	UserCache map[string]*model.User
)

func InitController() error {
	// Check if we want to activate DB
	if os.Getenv("DB_ACTIVE") == "true" {
		// Initialize the database
		err := model.InitDB()
		if err != nil {
			log.Errf("Error initializing database: %s", err)
			return err
		}

		// Load users from DB
		tCache, err := model.LoadUsers()
		if err != nil {
			return err
		}
		log.Logf("Loaded %d users from DB", len(*tCache))
		// Cache the users
		UserCache = *tCache
	}
	// Load users
	return nil
}
