package env

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

// Load initializes a .env environment
func Load() error {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		// Ignore missing env file in production
		if os.Getenv("PRODUCTION") == "" {
			return fmt.Errorf("environment initialization failed: %s", err)
		}
	}
	return nil
}
