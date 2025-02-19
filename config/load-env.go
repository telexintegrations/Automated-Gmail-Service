package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURI  string
}

var AppConfig Config

func LoadConfig() (configFile Config) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	AppConfig = Config{
		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURI:  os.Getenv("GOOGLE_REDIRECT_URI"),
	}
	return AppConfig
}
