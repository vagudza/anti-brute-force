package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vagudza/anti-brute-force/internal/config"
)

type Storage struct {
	pool *pgxpool.Pool
}

func NewStorage(ctx context.Context, cfg *config.PGConfig) (*Storage, error) {
	pool, err := pgxpool.New(ctx, buildDSN(cfg))
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	return &Storage{pool: pool}, nil
}

func buildDSN(cfg *config.PGConfig) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database, cfg.SSLMode)
}
