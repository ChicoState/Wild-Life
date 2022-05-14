package env

/*
* Part of the Ruby dotenv project
* Loads vars from a .env file
* Env files store anything that can change between different systems
 */

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Load initializes a .env environment
func Load() error {
	err := godotenv.Load()
	if err != nil {
		// Ignore missing env file in production
		if os.Getenv("PRODUCTION") == "" {
			return fmt.Errorf("environment initialization failed: %s", err)
		}
	}
	return nil
}
