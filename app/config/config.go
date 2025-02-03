package config

import (
	types "api/app/types/env"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var configInstance *types.EnvConfig
var configLoaded = false

func Set(cfg *types.EnvConfig) {
	configInstance = cfg
	configLoaded = true
}

func Get() *types.EnvConfig {
	if !configLoaded {
		Load()
	}
	return configInstance
}

func Load() *types.EnvConfig {
	if configLoaded {
		return configInstance
	}
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

	dbPortStr := os.Getenv("420_DB_PORT")
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil || dbPort == 0 {
		dbPort = 3306
	}
	log.Printf("Loaded DB_PORT: %d", dbPort)

	app := os.Getenv("420_APP")
	log.Printf("Loaded 420_APP: %s", app)

	db := os.Getenv("420_DB")
	log.Printf("Loaded 420_DB: %s", db)

	dbTable := os.Getenv("420_DB_TABLE")
	log.Printf("Loaded 420_DB_TABLE: %s", dbTable)

	configInstance = &types.EnvConfig{
		APIPort: apiPort,
		App:     app,
		DB:      db,
		DBPort:  dbPort,
		DBTable: dbTable,
	}

	configLoaded = true
	return configInstance
}

func Reload() {
	configLoaded = false
	Load()
}
