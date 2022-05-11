package controller

import "testing"

func TestInitController(t *testing.T) {
	err := InitController() //generally manipulated by a .env
	if err != nil {
		t.Error(err)
	}
}
