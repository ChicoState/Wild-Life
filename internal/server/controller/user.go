package controller

import (
	"net/http"
	"wildlife/internal/log"
	"wildlife/internal/server/model"
)

func TestDBAdd(w http.ResponseWriter, r *http.Request) {
	name := "test"
	email := "test@test.test"
	user := model.AddUser(name, email)
	err := user.Save()
	if err != nil { //couldn't save user to db
		log.Errf("Error saving user to db: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else { //passes all tests
		user_cache[user.ID] = &user
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("User added successfully"))
		if err != nil {
			log.Errf("Error writing response: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func TestDBRemove(w http.ResponseWriter, r *http.Request) {
	email := "test@test.test"
	//little slow but we can't assume cookies are enabled
	for _, user := range user_cache {
		if user.Email == email {
			err := user.Delete()
			if err != nil {
				log.Errf("Error deleting user from db: %s", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				delete(user_cache, user.ID)
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
