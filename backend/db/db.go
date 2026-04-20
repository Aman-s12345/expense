package db

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresDB struct {
	db *sqlx.DB
}

func NewPostgresConnection(host string) (*PostgresDB, error) {
	db, err := sqlx.Connect("postgres", host)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	db.SetMaxOpenConns(8)
    db.SetMaxIdleConns(4)
    db.SetConnMaxLifetime(5*time.Minute)
    db.SetConnMaxIdleTime(1*time.Minute)

	log.Info("Connected to PostgreSQL")
	if err := db.PingContext(context.Background()); err != nil {
        return nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
    }
	return &PostgresDB{db: db}, nil
}

func (p *PostgresDB) GetDB() *sqlx.DB {
	return p.db
}

func (p *PostgresDB) Close() error {
	return p.db.Close()
}

func (p *PostgresDB) Ping(ctx context.Context) error {
	return p.db.PingContext(ctx)
}

func (p *PostgresDB) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return p.db.BeginTxx(ctx, nil)
}

func (p *PostgresDB) PoolStats() map[string]interface{} {
    stats := p.db.Stats()
    return map[string]interface{}{
        "max_open_connections":   8, 
        "open_connections":       stats.OpenConnections,
        "in_use":                stats.InUse,
        "idle":                  stats.Idle,
        "wait_count":            stats.WaitCount,
        "wait_duration_ms":      stats.WaitDuration.Milliseconds(),
        "max_idle_closed":       stats.MaxIdleClosed,
        "max_lifetime_closed":   stats.MaxLifetimeClosed,
        "max_idle_time_closed":  stats.MaxIdleTimeClosed,
    }
}