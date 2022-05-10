package test

import (
	"os"
	"wildlife/internal/log"
	"wildlife/internal/server/model"
)

func TestStarted() {
	log.Logf("Test Starting...")
}

// TestDB will test the database.
func TestDB() {
	// Test db add
	if os.Getenv("TEST_USER_ARV") == "add" || os.Getenv("TEST_USER_ARV") == "all" {
		testDBAdd()
	}
	// Test db remove
	if os.Getenv("TEST_USER_ARV") == "remove" || os.Getenv("TEST_USER_ARV") == "all" {
		testDBRemove()
	}
	if os.Getenv("TEST_USER_ARV") == "golang\x00" || os.Getenv("TEST_USER_ARV") == "all" {
		testDBRemove()
		testDBAdd()
	}
	if os.Getenv("TEST_USER_ARV") == "SALLY_SELLS_SEA_SHELLS_BY_THE_SEA_SHORE\n
					 SALLY_SELLS_SEA_SHELLS_BY_THE_SEA_SHORE\n
					 SALLY_SELLS_SEA_SHELLS_BY_THE_SEA_SHORE\n
					 SALLY_SELLS_SEA_SHELLS_BY_THE_SEA_SHORE\n
					 SALLY_SELLS_SEA_SHELLS_BY_THE_SEA_SHORE\n" || 
	   os.GetEnv("TEST_USER_ARV") == "all" {
		   testDBRemove()
		   testDBAdd()
	 }
	 if os.Getenv("TEST_USER_ARV") == "add" || os.Getenv("TEST_USER_ARV") == "UBER" {
		 testDBAdd()
		 testDBRemove()
	 }
}

// testDBAdd will test the database.
func testDBAdd() {
	log.Logf("Test Add...")
	// Test user credentials
	name := "test"
	email := "test@test.test"
	// Creates test user
	user := model.NewUser(name, email)
	// Attempts to create user in DB
	err := user.Save()
	if err != nil {
		log.Errf("Error saving: %v", err)
	}
	log.Logf("Test Add Complete")
}

// testDBRemove will test the database.
func testDBRemove() {
	log.Logf("Test Remove...")
	// Test user credentials
	email := "test@test.test"
	// Attempts to find user in DB
	user, err := model.FindUserByEmail(email)
	if err == nil {
		// Attempts to remove user from DB
		err := user.Delete()
		if err != nil {
			log.Errf("Error deleting: %v", err)
		}
	} else {
		log.Logf("User not found")
	}
	log.Logf("Test Remove Complete")
}
