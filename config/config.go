package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

const (
	envFilename   = ".env"
	productionEnv = "production"
	demoEnv       = "demo"
)

func init() {
	env := os.Getenv("ENV")
	if env != productionEnv && env != demoEnv {
		err := godotenv.Load(envFilename)
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}

type Config struct {
	APIBaseURL string
	APIToken   string
	ServerPort string

	PostgresPort     string
	PostgresHost     string
	PostgresPassword string
	PostgresUser     string
	PostgresDatabase string

	//ShortSmaInterval int
	//LongSmaInterval  int
	//RSIInterval      int
}

func LoadConfig() *Config {
	return &Config{
		APIBaseURL: os.Getenv("TINKOFF_API_BASE_URL"),
		APIToken:   os.Getenv("TINKOFF_API_TOKEN"),
		ServerPort: os.Getenv("SERVER_PORT"),

		PostgresPort:     os.Getenv("POSTGRES_PORT"),
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresDatabase: os.Getenv("POSTGRES_DATABASE"),
	}
}
