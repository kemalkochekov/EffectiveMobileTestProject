package configs

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Port string
	}
	Postgres struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
	}
	ShutdownTimeout time.Duration
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, err
	}

	shutdownTimeout, err := time.ParseDuration(os.Getenv("SHUTDOWN_TIMEOUT"))
	if err != nil {
		return nil, err
	}

	c := &Config{
		Server: struct {
			Port string
		}{
			Port: os.Getenv("HTTP_PORT"),
		},
		Postgres: struct {
			Host     string
			Port     int
			User     string
			Password string
			DBName   string
		}{
			Host:     os.Getenv("DB_HOST"),
			Port:     dbPort,
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
		},
		ShutdownTimeout: shutdownTimeout,
	}

	return c, nil
}
