package config

type EnvVariable string

const (
	ServerPort EnvVariable = "SERVER_PORT"
	DBHost     EnvVariable = "DB_HOST"
	DBPort     EnvVariable = "DB_PORT"
	DBUser     EnvVariable = "DB_USER"
	DBPassword EnvVariable = "DB_PASSWORD"
	DBName     EnvVariable = "DB_NAME"
	DBSSLMode  EnvVariable = "DB_SSLMODE"
)
