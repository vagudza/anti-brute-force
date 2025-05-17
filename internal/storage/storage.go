package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vagudza/anti-brute-force/internal/config"
)

type Repository interface {
	AddSubnetToWhitelist(ctx context.Context, subnet string) error
	RemoveSubnetFromWhitelist(ctx context.Context, subnet string) error
	IsIPInWhitelist(ctx context.Context, ip string) (bool, error)

	AddSubnetToBlacklist(ctx context.Context, subnet string) error
	RemoveSubnetFromBlacklist(ctx context.Context, subnet string) error
	IsIPInBlacklist(ctx context.Context, ip string) (bool, error)
}

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

func (s *Storage) AddSubnetToWhitelist(ctx context.Context, subnet string) error {
	query := `
		INSERT INTO whitelist (subnet)
		VALUES ($1)
		ON CONFLICT (subnet) DO NOTHING
	`
	_, err := s.pool.Exec(ctx, query, subnet)
	return err
}

func (s *Storage) RemoveSubnetFromWhitelist(ctx context.Context, subnet string) error {
	query := `
		DELETE FROM whitelist
		WHERE subnet = $1
	`
	_, err := s.pool.Exec(ctx, query, subnet)
	return err
}

func (s *Storage) IsIPInWhitelist(ctx context.Context, ip string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM whitelist
			WHERE $1 <<= subnet::cidr
		)
	`
	var exists bool
	err := s.pool.QueryRow(ctx, query, ip).Scan(&exists)
	return exists, err
}

func (s *Storage) AddSubnetToBlacklist(ctx context.Context, subnet string) error {
	query := `
		INSERT INTO blacklist (subnet)
		VALUES ($1)
		ON CONFLICT (subnet) DO NOTHING
	`
	_, err := s.pool.Exec(ctx, query, subnet)
	return err
}

func (s *Storage) RemoveSubnetFromBlacklist(ctx context.Context, subnet string) error {
	query := `
		DELETE FROM blacklist
		WHERE subnet = $1
	`
	_, err := s.pool.Exec(ctx, query, subnet)
	return err
}

func (s *Storage) IsIPInBlacklist(ctx context.Context, ip string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM blacklist
			WHERE $1 <<= subnet::cidr
		)
	`
	var exists bool
	err := s.pool.QueryRow(ctx, query, ip).Scan(&exists)
	return exists, err
}
