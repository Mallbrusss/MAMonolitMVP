package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
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

	ShortSmaInterval int
	LongSmaInterval  int
	RSIInterval      int
}

func LoadConfig() *Config {
	shortSmaInterval, err := strconv.Atoi(os.Getenv("SHORTSMA_INTERVAL"))
	if err != nil {
		log.Fatal(err)
	}
	longSmaInterval, err := strconv.Atoi(os.Getenv("LONG_SMA_INTERVAL"))
	if err != nil {
		log.Fatal("Error loading LONG_SMA_INTERVAL")
	}
	rsiInterval, err := strconv.Atoi(os.Getenv("RSI_INTERVAL"))
	if err != nil {
		log.Fatal("Error loading RSI_INTERVAL")
	}

	return &Config{
		APIBaseURL: os.Getenv("TINKOFF_API_BASE_URL"),
		APIToken:   os.Getenv("TINKOFF_API_TOKEN"),
		ServerPort: os.Getenv("SERVER_PORT"),

		PostgresPort:     os.Getenv("POSTGRES_PORT"),
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresDatabase: os.Getenv("POSTGRES_DATABASE"),

		ShortSmaInterval: shortSmaInterval,
		LongSmaInterval:  longSmaInterval,
		RSIInterval:      rsiInterval,
	}
}
