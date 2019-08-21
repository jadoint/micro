package env

import (
	"os"

	"github.com/joho/godotenv"
)

// Load loads environment variables from .env
func Load() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}
	godotenv.Load(".env." + env)
	return nil
}
