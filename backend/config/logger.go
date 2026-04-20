package config

import (
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2/log"
)

type LoggerConfig struct {
	Level log.Level
}

func (c *Config) LoadLoggerConfig() {
	level := log.LevelInfo

	if levelStr := os.Getenv("LOG_LEVEL"); levelStr != "" {
		if lvl, err := strconv.Atoi(levelStr); err == nil {
			level = log.Level(lvl)
		}
	}

	log.SetLevel(level)
	c.Logger.Level = level
}
