package config

import (
	"420-api/app/types"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var configInstance *types.Config

func Set(cfg *types.Config) {
	configInstance = cfg
}

func Get() *types.Config {
	return configInstance
}

func Load() *types.Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	apiPortStr := os.Getenv("420_API_PORT")
	apiPort, err := strconv.Atoi(apiPortStr)
	if err != nil || apiPort == 0 {
		apiPort = 8080
	}
	log.Printf("Loaded API_PORT: %d", apiPort)

	app := os.Getenv("420_APP")
	log.Printf("Loaded 420_APP: %s", app)

	return &types.Config{
		APIPort: apiPort,
		AllowedOrigin: types.AllowedOrigin{
			App: app,
		},
	}
}
