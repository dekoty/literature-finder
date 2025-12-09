package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	APIKey        string
	DatabaseURL   string
	ServerAddress string
}

func Load() (*Config, error) {

	if err := godotenv.Load(); err != nil {
		log.Println("Info: файл .env не найден, используются системные переменные")
	}

	dsn := os.Getenv("DATASOURCENAME")
	if dsn == "" {
		return nil, fmt.Errorf("DATASOURCENAME не установлен")
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Println("Warning: API_KEY не установлен, функционал Google Books может быть ограничен")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	return &Config{
		APIKey:        apiKey,
		DatabaseURL:   dsn,
		ServerAddress: port,
	}, nil
}
