package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	DatabaseURL string
	ServerHost  string
	ServerPort  string
	ServerURL   string
	JwtSecret   string
	JwtExpire   string
}

var AppConfig = config{}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Cannot read .env file")
	}

	AppConfig.JwtSecret = tryGetEnv("JWT_SECRET", "30bd725869cb74ba")
	AppConfig.JwtExpire = tryGetEnv("JWT_EXPIRE", "24h")
	AppConfig.DatabaseURL = tryGetEnv("DATABASE_URL", "localhost")
	AppConfig.ServerHost = tryGetEnv("SERVER_HOST", "localhost")
	AppConfig.ServerPort = tryGetEnv("SERVER_PORT", "6900")
	AppConfig.ServerURL = fmt.Sprintf("http://%s:%s", AppConfig.ServerHost, AppConfig.ServerPort)

	log.Printf("Config successfully loaded: %+v", AppConfig)
}

// Gets an environment variable by key if exists, otherwise returns a default value
func tryGetEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
