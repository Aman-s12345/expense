package config

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

type Config struct {
	Server     ServerConfig
	DB         DBConfig
	Logger     LoggerConfig	
}

func NewConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Info("No .env file found, using environment variables")
	}

	cfg := &Config{}
	cfg.LoadServerConfig()
	cfg.LoadDBConfig()


	return cfg
}
