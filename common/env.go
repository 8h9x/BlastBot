package common

import (
	"github.com/joho/godotenv"
	"os"
)

// only loads .env file if not prod
func LoadEnv() error {
	_, isProd := os.LookupEnv("PROD")

	if !isProd {
		err := godotenv.Load(".env")
		if err != nil {
			return err
		}
	}

	return nil
}
