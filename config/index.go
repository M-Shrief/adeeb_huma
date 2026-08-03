package config

import (
	"crypto/rsa"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/golang-jwt/jwt"
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

// Cache
var (
	Cache_HOST     string
	Cache_PORT     int
	Cache_Username string
	Cache_PASSWORD string
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

	Cache_HOST = envs["Cache_HOST"]
	Cache_PORT, _ = strconv.Atoi(envs["Cache_PORT"])
	Cache_Username = envs["CacCachCache_Usernamee_Usernamehe_HOST"]
	Cache_PASSWORD = envs["Cache_PASSWORD"]

	assignJWTKeys()
}

// JWT
var (
	JWT_PRIVATE_KEY  *rsa.PrivateKey
	JWT_PUBLIC_KEY   *rsa.PublicKey
	JWT_PRIVATE_PATH string
	JWT_PUBLIC_PATH  string
)

func assignJWTKeys() {

	JWT_PRIVATE_PATH = envs["JWT_PRIVATE_PATH"]
	JWT_PUBLIC_PATH = envs["JWT_PUBLIC_PATH"]

	jwt_Private_File, err := os.ReadFile(JWT_PRIVATE_PATH)
	if err != nil {
		fmt.Printf("couldn't parse private JWT secret file: %v\n", err)
		os.Exit(1)
	}

	JWT_PRIVATE_KEY, err = jwt.ParseRSAPrivateKeyFromPEM(jwt_Private_File)
	if err != nil {
		fmt.Printf("couldn't parse private JWT secret: %v\n", err)
		os.Exit(1)
	}

	// JWT
	jwt_Public_File, err := os.ReadFile(JWT_PUBLIC_PATH)
	if err != nil {
		fmt.Printf("couldn't read public JWT secret file: %v\n", err)
		os.Exit(1)
	}

	JWT_PUBLIC_KEY, err = jwt.ParseRSAPublicKeyFromPEM(jwt_Public_File)
	if err != nil {
		fmt.Printf("couldn't parse public JWT secret: %v\n", err)
		os.Exit(1)
	}
}
