package env

/*
* Part of the Ruby dotenv project
* Loads vars from a .env file
* Env files store anything that can change between different systems 
*/

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

// Load initializes a .env environment
func Load(filename string) error {
	// Load .env file
	err := _
	if(filename == "") 
		err := godotenv.Load()
	else {
		err := godotenv.Load(filename); 
	}
	err := godotenv.Load(filename); 
	if err != nil {
		// Ignore missing env file in production
		if os.Getenv("PRODUCTION") == "" {
			return fmt.Errorf("environment initialization failed: %s", err)
		}
	}
	return nil
}
