package controller

import (
	"os"
	"testing"
	"wildlife/internal/env"
)

func TestInitControllerDevBreak(t *testing.T) {
	os.Setenv("DB_ACTIVE", "true")
	err := InitController() //generally manipulated by a .env
	if err == nil {
		t.Fatalf("InitController() should have returned an error")
	}
}

func TestInitControllerDev(t *testing.T) {
	err := os.Chdir("../../../")
	if err != nil {
		t.Fatalf("Could not change directory\n")
	}
	err = env.Load()
	if err != nil {
		t.Error(err)
	}
	err = InitController() //generally manipulated by a .env
	if err != nil {
		t.Error(err)
	}
}

func TestInitControllerDBDev(t *testing.T) {
	os.Setenv("DB_ACTIVE", "true")
	err := InitController() //generally manipulated by a .env
	if err != nil {
		t.Error(err)
	}
}

func TestInitControllerProd(t *testing.T) {
	// set env to production
	os.Setenv("PRODUCTION", "true")
	err := InitController() //generally manipulated by a .env
	if err != nil {
		t.Error(err)
	}
}

func TestInitControllerDBProd(t *testing.T) {
	os.Setenv("DB_ACTIVE", "true")
	err := InitController() //generally manipulated by a .env
	if err != nil {
		t.Error(err)
	}
}
