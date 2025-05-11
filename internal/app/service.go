package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/vagudza/anti-brute-force/internal/bucket"
)

var (
	ErrEmptyLogin    = errors.New("empty login")
	ErrEmptyPassword = errors.New("empty password")
	ErrEmptyIP       = errors.New("empty IP")
)

type LimiterService interface {
	CheckAuth(ctx context.Context, login, password, ip string) (bool, error)
	ResetBucket(ctx context.Context, login, ip string) error
}

// Service представляет основную бизнес-логику приложения
type Service struct {
	loginBuckets    bucket.Limiter
	passwordBuckets bucket.Limiter
	ipBuckets       bucket.Limiter
	// whitelist       iplist.IPList
	// blacklist       iplist.IPList
}

// NewService создает новый экземпляр сервиса
func NewService(
	loginBuckets bucket.Limiter,
	passwordBuckets bucket.Limiter,
	ipBuckets bucket.Limiter,
	// whitelist iplist.IPList,
	// blacklist iplist.IPList,
) *Service {
	return &Service{
		loginBuckets:    loginBuckets,
		passwordBuckets: passwordBuckets,
		ipBuckets:       ipBuckets,
		// whitelist:       whitelist,
		// blacklist:       blacklist,
	}
}

// CheckAuth проверяет попытку авторизации
func (s *Service) CheckAuth(ctx context.Context, login, password, ip string) (bool, error) {
	if err := validateCheckAuth(login, password, ip); err != nil {
		return false, err
	}

	// Парсим IP
	// ip := net.ParseIP(req.IP)
	// if ip == nil {
	// 	return nil, fmt.Errorf("invalid IP address: %s", req.IP)
	// }

	// // Проверяем white/black листы
	// inWhitelist, err := s.whitelist.Contains(ctx, ip)
	// if err != nil {
	// 	return nil, fmt.Errorf("whitelist check error: %w", err)
	// }
	//
	// if inWhitelist {
	// 	return &models.AuthResponse{OK: true}, nil
	// }
	//
	// inBlacklist, err := s.blacklist.Contains(ctx, ip)
	// if err != nil {
	// 	return nil, fmt.Errorf("blacklist check error: %w", err)
	// }
	//
	// if inBlacklist {
	// 	return &models.AuthResponse{OK: false}, nil
	// }

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

// ResetBucket сбрасывает bucket для указанного логина и IP
func (s *Service) ResetBucket(ctx context.Context, login, ip string) error {
	if login != "" {
		if err := s.loginBuckets.Reset(ctx, login); err != nil {
			return fmt.Errorf("login bucket reset error: %w", err)
		}
	}

	if ip != "" {
		if err := s.ipBuckets.Reset(ctx, ip); err != nil {
			return fmt.Errorf("IP bucket reset error: %w", err)
		}
	}

	return nil
}

//
// // AddToBlacklist добавляет подсеть в черный список
// func (s *Service) AddToBlacklist(ctx context.Context, subnet *net.IPNet) error {
// 	return s.blacklist.Add(ctx, subnet)
// }
//
// // RemoveFromBlacklist удаляет подсеть из черного списка
// func (s *Service) RemoveFromBlacklist(ctx context.Context, subnet *net.IPNet) error {
// 	return s.blacklist.Remove(ctx, subnet)
// }
//
// // AddToWhitelist добавляет подсеть в белый список
// func (s *Service) AddToWhitelist(ctx context.Context, subnet *net.IPNet) error {
// 	return s.whitelist.Add(ctx, subnet)
// }
//
// // RemoveFromWhitelist удаляет подсеть из белого списка
// func (s *Service) RemoveFromWhitelist(ctx context.Context, subnet *net.IPNet) error {
// 	return s.whitelist.Remove(ctx, subnet)
// }

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

	return nil
}
