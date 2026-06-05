package config

import (
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	DatabaseURL  string
	JWTSecret    string
	RefreshSecret string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Предупреждение: Файл .env не найден, используются системные переменные окружения")
	}

	dbUser := getEnv("DB_USER", "postgres")
	dbPass := getEnv("DB_PASSWORD", "")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbName := getEnv("DB_NAME", "ozinse")

	// Красиво собираем стандартную строку подключения (DSN) для lib/pq
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", 
		dbUser, dbPass, dbHost, dbPort, dbName,
	)

	return &Config{
		Port:          getEnv("PORT", "8080"),
		DatabaseURL:   dsn,
		JWTSecret:     getEnv("JWT_SECRET", "default_access_secret"),
		RefreshSecret: getEnv("REFRESH_SECRET", "default_refresh_secret"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}