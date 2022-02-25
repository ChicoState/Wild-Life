package controller

import (
	"wildlife/internal/log"
	model "wildlife/internal/server/model"
)

var (
	user_cache map[string]*model.User
)

func InitController() {
	t_cache, err := model.LoadUsers()
	if err != nil {
		log.Errf("Error loading users: %s", err)
		panic(err)
	}
	user_cache = *t_cache
}
