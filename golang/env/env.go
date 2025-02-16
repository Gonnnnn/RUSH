// Helper package to handle environment variables.
package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/ridge/must/v2"
)

// The function that should be called first in the main function.
// It loads the environment variables from the file.
func Load(fileNameKey string) {
	envFile := GetOptionalStringVariable(fileNameKey, "")
	if envFile != "" {
		must.OK(godotenv.Load(envFile))
	}
}

func GetRequiredStringVariable(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Fatalln("Environment variable (" + name + ") is required but does not exist.")
	}
	return value
}

func GetOptionalStringVariable(name string, defaultValue string) string {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}
	return value
}
