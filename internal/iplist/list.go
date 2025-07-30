package iplist

import (
	"context"
	"errors"
	"fmt"
	"net/netip"

	"github.com/vagudza/anti-brute-force/internal/storage"
)

var (
	ErrInvalidIP     = errors.New("invalid IP address format")
	ErrInvalidSubnet = errors.New("invalid subnet CIDR format")
)

type ServiceClient interface {
	AddToWhitelist(ctx context.Context, subnet string) error
	RemoveFromWhitelist(ctx context.Context, subnet string) error
	ContainsInWhitelist(ctx context.Context, ip string) (bool, error)
	GetWhitelist(ctx context.Context) ([]string, error)
	ClearWhitelist(ctx context.Context) error

	AddToBlacklist(ctx context.Context, subnet string) error
	RemoveFromBlacklist(ctx context.Context, subnet string) error
	ContainsInBlacklist(ctx context.Context, ip string) (bool, error)
	GetBlacklist(ctx context.Context) ([]string, error)
	ClearBlackList(ctx context.Context) error
}

type Service struct {
	repo storage.Repository
}

func NewService(repo storage.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) AddToWhitelist(ctx context.Context, subnet string) error {
	if _, err := netip.ParsePrefix(subnet); err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidSubnet, subnet)
	}

	return s.repo.AddSubnetToWhitelist(ctx, subnet)
}

func (s *Service) RemoveFromWhitelist(ctx context.Context, subnet string) error {
	if _, err := netip.ParsePrefix(subnet); err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidSubnet, subnet)
	}

	return s.repo.RemoveSubnetFromWhitelist(ctx, subnet)
}

func (s *Service) ContainsInWhitelist(ctx context.Context, ip string) (bool, error) {
	if _, err := netip.ParseAddr(ip); err != nil {
		return false, ErrInvalidIP
	}

	return s.repo.IsIPInWhitelist(ctx, ip)
}

func (s *Service) AddToBlacklist(ctx context.Context, subnet string) error {
	if _, err := netip.ParsePrefix(subnet); err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidSubnet, subnet)
	}

	return s.repo.AddSubnetToBlacklist(ctx, subnet)
}

func (s *Service) RemoveFromBlacklist(ctx context.Context, subnet string) error {
	if _, err := netip.ParsePrefix(subnet); err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidSubnet, subnet)
	}

	return s.repo.RemoveSubnetFromBlacklist(ctx, subnet)
}

func (s *Service) ContainsInBlacklist(ctx context.Context, ip string) (bool, error) {
	if _, err := netip.ParseAddr(ip); err != nil {
		return false, ErrInvalidIP
	}

	return s.repo.IsIPInBlacklist(ctx, ip)
}

func (s *Service) GetWhitelist(ctx context.Context) ([]string, error) {
	return s.repo.GetWhitelist(ctx)
}

func (s *Service) GetBlacklist(ctx context.Context) ([]string, error) {
	return s.repo.GetBlacklist(ctx)
}

func (s *Service) ClearWhitelist(ctx context.Context) error {
	return s.repo.ClearWhiteList(ctx)
}

func (s *Service) ClearBlackList(ctx context.Context) error {
	return s.repo.ClearBlackList(ctx)
}
