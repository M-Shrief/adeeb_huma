package config

import (
	"log"

	"github.com/joho/godotenv"
)

var envs map[string]string

func LoadENV() {
	var err error
	envs, err = godotenv.Read(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	assignValues()
}

// General
var (
	APP_ENV string
	HOST    string
	PORT    string
)

// Database
var (
	DB_HOST     string
	DB_PORT     string
	DB_NAME     string
	DB_USER     string
	DB_PASSWORD string
)

func assignValues() {
	// General
	APP_ENV = envs["APP_ENV"]
	HOST = envs["HOST"]
	PORT = envs["PORT"]

	// Database
	DB_HOST = envs["DB_HOST"]
	DB_PORT = envs["DB_PORT"]
	DB_USER = envs["DB_USER"]
	DB_NAME = envs["DB_NAME"]
	DB_PASSWORD = envs["DB_PASSWORD"]
}
