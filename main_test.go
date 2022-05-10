package main

import (
	"testing"
	"wildlife/internal/env"
)

//Test Enviornment
// t parameter Fatalf called if func ret no error or non empty str
func TestEnv(t *testing.T) {
	err := env.Load()
	if err != nil {
		t.Fatalf("Environment initialization failed: %s", err)
	}
}

//Test controller
//See more tests of the database in internal/server/controller/user.go
// func TestController(t *testing.T) {
// 	os.Setenv("DB_ACTIVE", "false")
// 	//pretty much starts up the DB
// 	err := controller.InitController()
// 	if err != nil {
// 		t.Fatalf("Database is not active. D00d.\n")
// 	}
// }

//Test the web server
//See more tests in internal/server/server.go
// func TestServerPort(t *testing.T) {
// 	os.Setenv("PORT", "42069")
// 	err := server.Start()
// 	if err != nil {
// 		t.Fatalf("Teehee\n")
// 	}
// }
