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
	GetWhitelist(ctx context.Context) ([]string, error)

	AddSubnetToBlacklist(ctx context.Context, subnet string) error
	RemoveSubnetFromBlacklist(ctx context.Context, subnet string) error
	IsIPInBlacklist(ctx context.Context, ip string) (bool, error)
	GetBlacklist(ctx context.Context) ([]string, error)
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

func (s *Storage) GetWhitelist(ctx context.Context) ([]string, error) {
	query := `
		SELECT subnet::text
		FROM whitelist
		ORDER BY created_at DESC
	`
	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query whitelist: %w", err)
	}
	defer rows.Close()

	var subnets []string
	for rows.Next() {
		var subnet string
		if err := rows.Scan(&subnet); err != nil {
			return nil, fmt.Errorf("failed to scan whitelist subnet: %w", err)
		}
		subnets = append(subnets, subnet)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating whitelist rows: %w", err)
	}

	return subnets, nil
}

func (s *Storage) GetBlacklist(ctx context.Context) ([]string, error) {
	query := `
		SELECT subnet::text
		FROM blacklist
		ORDER BY created_at DESC
	`
	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query blacklist: %w", err)
	}
	defer rows.Close()

	var subnets []string
	for rows.Next() {
		var subnet string
		if err := rows.Scan(&subnet); err != nil {
			return nil, fmt.Errorf("failed to scan blacklist subnet: %w", err)
		}
		subnets = append(subnets, subnet)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating blacklist rows: %w", err)
	}

	return subnets, nil
}
