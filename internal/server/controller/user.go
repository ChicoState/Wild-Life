package controller

import (
	"net/http"
	"wildlife/internal/log"
	"wildlife/internal/server/model"
)

// Tests adding a user to the DB
func TestDBAdd(w http.ResponseWriter, r *http.Request) {
	// Test user credentials
	name := "test"
	email := "test@test.test"
	// Creates test user
	user := model.NewUser(name, email)
	// Attempts to create user in DB
	err := user.Save()
	if err != nil { //couldn't save user to db
		log.Errf("Error saving user to db: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else { //passes all tests
		// sets user in cache
		UserCache[user.ID] = user
		// Sets status code to OK
		w.WriteHeader(http.StatusOK)
		// prints basic message to browser
		_, err := w.Write([]byte("User added successfully"))
		if err != nil {
			log.Errf("Error writing response: %s", err)
			// prints basic message to browser
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// Tests deleting a user from the DB
func TestDBRemove(w http.ResponseWriter, r *http.Request) {
	// The test user email
	email := "test@test.test"
	// little slow but we can't assume cookies are enabled
	// and we don't know the session id
	for _, user := range UserCache {
		// if the user email matches the test email
		if user.Email == email {
			// delete the user from the cache
			err := user.Delete()
			if err != nil {
				log.Errf("Error deleting user from db: %s", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				// delete the user from the cache
				delete(UserCache, user.ID)
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte("User deleted successfully"))
				if err != nil {
					log.Errf("Error writing response: %s", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
			return
		}
	}
	log.Errf("User not found")
	http.Error(w, "User not found", http.StatusNotFound)
}

// Tests showing database users
func TestDBView(w http.ResponseWriter, r *http.Request) {
	// Creates html string from cache
	response := "<html><body><h1>Cache Users</h1><ul>"
	for _, user := range UserCache {
		response += "<li>" + user.Email + "</li>"
	}
	response += "</ul></body></html>"
	w.WriteHeader(http.StatusOK)
	// prints basic message to browser
	_, err := w.Write([]byte(response))
	if err != nil {
		log.Errf("Error writing response: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
