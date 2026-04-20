package config

type DBConfig struct {
	Host string
	Name string
}

func (c *Config) LoadDBConfig() {
	c.DB.Host = getEnv("DB_HOST", "postgresql://postgres:password@localhost:5432/wealth_be?sslmode=disable")
	c.DB.Name = getEnv("DB_NAME", "")
}
