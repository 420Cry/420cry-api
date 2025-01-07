package config

import (
	"420-api/app/types"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func Load() *types.Config {
	godotenv.Load()

	apiPortStr := os.Getenv("420_API_PORT")
	apiPort, err := strconv.Atoi(apiPortStr)
	if err != nil || apiPort == 0 {
		apiPort = 8080
	}
	log.Printf("Loaded API_PORT: %d", apiPort)

	app := os.Getenv("420_APP")
	log.Printf("Loaded 420_APP: %s", app)

	return &types.Config{
		API_PORT: apiPort,
		ALLOWED_ORIGIN: types.AllowedOrigin{
			APP: app,
		},
	}
}
