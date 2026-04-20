package providers

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"

	"github.com/Aman-s12345/expense/backend/config"
	"github.com/Aman-s12345/expense/backend/db"
)

type Provider struct {
	Config *config.Config
	DB     *db.PostgresDB
	S      *Services
}

func NewProvider(cfg *config.Config) (*Provider, error) {
	postgresDB, err := db.NewPostgresConnection(cfg.DB.Host)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize PostgreSQL: %w", err)
	}

	if err := postgresDB.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	log.Info("Connected to PostgreSQL successfully")

	provider := &Provider{
		Config: cfg,
		DB:     postgresDB,
	}

	provider.S = NewServices(provider)
	log.Info("All services initialized successfully")

	return provider, nil
}

func (p *Provider) Cleanup() {
	if p.DB != nil {
		if err := p.DB.Close(); err != nil {
			log.Errorw("Failed to close PostgreSQL connection", "error", err)
		}
	}
}

type keyType struct {
	key string
}

var providerKey = keyType{"providers"}

func Handle(p *Provider) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Locals(providerKey, p)
		return c.Next()
	}
}

func GetProviders(c *fiber.Ctx) *Provider {
	return c.Locals(providerKey).(*Provider)
}