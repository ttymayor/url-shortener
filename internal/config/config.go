package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Redis    RedisConfig
	Postgres PostgresConfig
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using defaults or environment variables")
	}

	redisPort, _ := strconv.Atoi(getEnv("REDIS_PORT", "6379"))
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))

	pgPort, _ := strconv.Atoi(getEnv("POSTGRES_PORT", "5432"))

	return &Config{
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     redisPort,
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       redisDB,
		},
		Postgres: PostgresConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     pgPort,
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD", "postgres"),
			DBName:   getEnv("POSTGRES_DB", "url_shortener"),
			SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
