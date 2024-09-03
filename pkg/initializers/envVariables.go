package initializers

import (
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func LoadEnvVariables() {
	if _, runningInDocker := os.LookupEnv("RUNNING_IN_DOCKER"); !runningInDocker {
		err := godotenv.Load()
		if err != nil {
			log.Warn("Error loading .env file, proceeding with environment variables.")
		}
	} else {
		log.Info("Running in Docker, using environment variables provided by Docker.")
	}
}