package config

import (
	types "cry-api/app/types/env"
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

	dbPortStr := os.Getenv("DB_PORT")
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil || dbPort == 0 {
		dbPort = 3306
	}
	log.Printf("Loaded DB_PORT: %d", dbPort)

	app := os.Getenv("420_APP")
	log.Printf("Loaded 420_APP: %s", app)

	dbHost := os.Getenv("DB_HOST")
	log.Printf("Loaded DB_HOST: %s", dbHost)

	db := os.Getenv("DB_DATABASE")
	log.Printf("Loaded DB_DATABASE: %s", db)

	mysqlUser := os.Getenv("DB_USERNAME")
	log.Printf("Loaded DB_USERNAME: %s", mysqlUser)

	dbPassword := os.Getenv("DB_PASSWORD")
	log.Printf("Loaded DB_PASSWORD: %s", dbPassword)

	configInstance = &types.EnvConfig{
		APIPort:    apiPort,
		App:        app,
		DBHost:     dbHost,
		DBPort:     dbPort,
		DBDatabase: db,
		DBUserName: mysqlUser,
		DBPassword: dbPassword,
	}

	configLoaded = true
	return configInstance
}

func Reload() {
	configLoaded = false
	Load()
}
