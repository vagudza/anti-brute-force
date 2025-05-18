package app

import (
	"context"
	"errors"
	"fmt"
	"net"

	"go.uber.org/zap"

	"github.com/vagudza/anti-brute-force/internal/bucket"
	"github.com/vagudza/anti-brute-force/internal/iplist"
)

var (
	ErrEmptyLogin    = errors.New("empty login")
	ErrEmptyPassword = errors.New("empty password")
	ErrEmptyIP       = errors.New("empty IP")
	ErrInvalidIP     = errors.New("invalid IP address")
)

type LimiterService interface {
	CheckAuth(ctx context.Context, login, password, ip string) (bool, error)
	ResetBucket(ctx context.Context, login, ip string) error

	AddToWhitelist(ctx context.Context, subnet string) error
	RemoveFromWhitelist(ctx context.Context, subnet string) error
	AddToBlacklist(ctx context.Context, subnet string) error
	RemoveFromBlacklist(ctx context.Context, subnet string) error
}

type Service struct {
	logger          *zap.Logger
	loginBuckets    bucket.Limiter
	passwordBuckets bucket.Limiter
	ipBuckets       bucket.Limiter
	ipListService   iplist.ServiceClient
}

func NewService(
	logger *zap.Logger,
	loginBuckets bucket.Limiter,
	passwordBuckets bucket.Limiter,
	ipBuckets bucket.Limiter,
	ipListService iplist.ServiceClient,
) *Service {
	return &Service{
		logger:          logger,
		loginBuckets:    loginBuckets,
		passwordBuckets: passwordBuckets,
		ipBuckets:       ipBuckets,
		ipListService:   ipListService,
	}
}

func (s *Service) CheckAuth(ctx context.Context, login, password, ip string) (bool, error) {
	if err := validateCheckAuth(login, password, ip); err != nil {
		return false, err
	}

	inWhitelist, err := s.ipListService.ContainsInWhitelist(ctx, ip)
	if err != nil {
		return false, fmt.Errorf("whitelist check error: %w", err)
	}

	if inWhitelist {
		return true, nil
	}

	inBlacklist, err := s.ipListService.ContainsInBlacklist(ctx, ip)
	if err != nil {
		return false, fmt.Errorf("blacklist check error: %w", err)
	}

	if inBlacklist {
		return false, nil
	}

	loginAllowed, err := s.loginBuckets.Allow(ctx, login)
	if err != nil {
		return false, fmt.Errorf("login check error: %w", err)
	}

	if !loginAllowed {
		return false, nil
	}

	passwordAllowed, err := s.passwordBuckets.Allow(ctx, password)
	if err != nil {
		return false, fmt.Errorf("password check error: %w", err)
	}

	if !passwordAllowed {
		return false, nil
	}

	ipAllowed, err := s.ipBuckets.Allow(ctx, ip)
	if err != nil {
		return false, fmt.Errorf("IP check error: %w", err)
	}

	if !ipAllowed {
		return false, nil
	}

	return true, nil
}

func (s *Service) ResetBucket(ctx context.Context, login, ip string) error {
	if login != "" {
		if err := s.loginBuckets.Reset(ctx, login); err != nil {
			return fmt.Errorf("login bucket reset error: %w", err)
		}
		s.logger.Info("login bucket reset", zap.String("login", login))
	}

	if ip != "" {
		if err := s.ipBuckets.Reset(ctx, ip); err != nil {
			return fmt.Errorf("IP bucket reset error: %w", err)
		}
		s.logger.Info("IP bucket reset", zap.String("IP", ip))
	}

	return nil
}

func (s *Service) AddToWhitelist(ctx context.Context, subnet string) error {
	err := s.ipListService.AddToWhitelist(ctx, subnet)
	if err != nil {
		return fmt.Errorf("whitelist add error: %w", err)
	}

	s.logger.Info("whitelist added", zap.String("subnet", subnet))
	return nil
}

func (s *Service) RemoveFromWhitelist(ctx context.Context, subnet string) error {
	err := s.ipListService.RemoveFromWhitelist(ctx, subnet)
	if err != nil {
		return fmt.Errorf("whitelist remove error: %w", err)
	}

	s.logger.Info("whitelist removed", zap.String("subnet", subnet))
	return nil
}

func (s *Service) AddToBlacklist(ctx context.Context, subnet string) error {
	err := s.ipListService.AddToBlacklist(ctx, subnet)
	if err != nil {
		return fmt.Errorf("blacklist add error: %w", err)
	}

	s.logger.Info("blacklist added", zap.String("subnet", subnet))
	return nil
}

func (s *Service) RemoveFromBlacklist(ctx context.Context, subnet string) error {
	err := s.ipListService.RemoveFromBlacklist(ctx, subnet)
	if err != nil {
		return fmt.Errorf("blacklist remove error: %w", err)
	}

	s.logger.Info("blacklist removed", zap.String("subnet", subnet))
	return nil
}

func validateCheckAuth(login, password, ip string) error {
	if login == "" {
		return ErrEmptyLogin
	}

	if password == "" {
		return ErrEmptyPassword
	}

	if ip == "" {
		return ErrEmptyIP
	}

	if net.ParseIP(ip) == nil {
		return ErrInvalidIP
	}

	return nil
}
