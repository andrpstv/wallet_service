package config

import (
	"fmt"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func Load() (*Config, error) {
	LoadEnv()

	cfg := &Config{
		Server: ServerConfig{
			Port: GetEnv(string(ServerPort), "8080"),
		},
		Database: DatabaseConfig{
			Host:     GetEnv(string(DBHost), "localhost"),
			Port:     GetEnv(string(DBPort), "5432"),
			User:     GetEnv(string(DBUser), "postgres"),
			Password: GetEnv(string(DBPassword), "password"),
			DBName:   GetEnv(string(DBName), "wallet_db"),
			SSLMode:  GetEnv(string(DBSSLMode), "disable"),
		},
	}

	return cfg, nil
}

func GetDBConnString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		GetEnv(string(DBHost), "localhost"),
		GetEnv(string(DBPort), "5432"),
		GetEnv(string(DBUser), "postgres"),
		GetEnv(string(DBPassword), "password"),
		GetEnv(string(DBName), "wallet_db"),
		GetEnv(string(DBSSLMode), "disable"))
}
