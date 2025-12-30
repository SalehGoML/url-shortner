package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl     string
	JWTSecret string
	Port      string
}

func LoadConfig() Config {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("no env file found, reading from env variables")
	}

	cfg := Config{
		DBUrl:     os.Getenv("DB_URL"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		Port:      os.Getenv("PORT"),
	}

	if cfg.DBUrl == "" || cfg.JWTSecret == "" || cfg.Port == "" {
		log.Fatal("Missing required env variables")
	}

	return cfg
}
