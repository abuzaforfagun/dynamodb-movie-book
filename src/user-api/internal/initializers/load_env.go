package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvVariables(filePath string) {
	err := godotenv.Load(filePath)
	if err != nil {
		log.Fatalf("Error loading '%s' file", filePath)
	}
}
