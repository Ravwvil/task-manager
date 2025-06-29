package config

import (
	"os"
	"time"
)

type Config struct {
	Host                      string
	Port                      string
	ReadTimeout               time.Duration
	WriteTimeout              time.Duration
	IdleTimeout               time.Duration
	GracefulShutdownTimeout   time.Duration
}

func New() *Config {
	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		Host:                      host,
		Port:                      port,
		ReadTimeout:               15 * time.Second,
		WriteTimeout:              15 * time.Second,
		IdleTimeout:               60 * time.Second,
		GracefulShutdownTimeout:   30 * time.Second,
	}
}
