package model

import (
	"os"
	"testing"
	"wildlife/internal/env"
)

var (
	testUser *User
)

func TestInitModelFailure(t *testing.T) {
	// set env to production
	err := InitDB() //generally manipulated by a .env
	if err == nil {
		t.Fatalf("InitModel() should have returned an error")
	}
	_, err = LoadUsers()
	if err == nil {
		t.Fatalf("LoadUsers() should have returned an error")
	}
	user := NewUser("test", "test")
	if user != nil {
		t.Fatalf("NewUser() should have returned nil")
	}
	user = &User{
		ID:    "test",
		Name:  "test",
		Email: "test",
	}
	err = user.Create()
	if err == nil {
		t.Fatalf("User.Create() should have returned an error")
	}
	err = user.Save()
	if err == nil {
		t.Fatalf("User.Save() should have returned an error")
	}
	err = user.Delete()
	if err == nil {
		t.Fatalf("User.Delete() should have returned an error")
	}
	_, err = FindUserByEmail("test")
	if err == nil {
		t.Fatalf("FindUserByEmail() should have returned nil")
	}
}

func TestInitModelSuccess(t *testing.T) {
	// set env to production
	err := os.Chdir("../../../")
	if err != nil {
		t.Fatalf("Could not change directory\n")
	}
	err = env.Load()
	if err != nil {
		t.Error(err)
	}
	err = InitDB() //generally manipulated by a .env
	if err != nil {
		t.Error(err)
	}
	testUser := NewUser("test", "test")
	if testUser == nil {
		t.Fatalf("NewUser() should have returned a user")
	}
	err = testUser.Create()
	if err != nil {
		t.Fatalf("User.Create() should have returned nil: %s", err)
	}
	err = testUser.Save()
	if err != nil {
		t.Fatalf("User.Save() should have returned nil: %s", err)
	}
	testUsers, err := LoadUsers()
	if err != nil {
		t.Error(err)
	}
	if len(*testUsers) == 0 {
		t.Fatalf("LoadUsers() should have returned at least one user")
	}
	temp, err := FindUserByEmail("test")
	if err != nil {
		t.Fatalf("FindUserByEmail() should have returned a user")
	}
	if testUser.ID != temp.ID {
		t.Fatalf("FindUserByEmail() should have returned the correct user: %s", testUser.ID)
	}
	err = testUser.Delete()
	if err != nil {
		t.Fatalf("User.Delete() should have returned nil")
	}
}
