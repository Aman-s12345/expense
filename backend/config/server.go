package config

import "os"

type ServerConfig struct {
	Host        string
	Port        string
	FrontendURL string
	BackendURL  string
	Name        string
}

type DeploymentConfig struct {
	Environment string
	Name        string
}

func (c *Config) LoadServerConfig() {
	c.Server.Host = getEnv("SERVER_HOST", "localhost")
	c.Server.Port = getEnv("SERVER_PORT", "3000")
	c.Server.FrontendURL = getEnv("FRONTEND_URL", "http://localhost:5173")
	c.Server.BackendURL = getEnv("BACKEND_URL", "http://localhost:3000")
	c.Server.Name = getEnv("SERVER_NAME", "Expense Backend")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
