package app

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"anti-brutforce/internal/bucket"
	"anti-brutforce/internal/entity"
)

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

func (s *Service) Register(srv *grpc.Server) {
	// statementv1.RegisterStatementServiceServer(srv, s)
}

// CheckAuth проверяет попытку авторизации
func (s *Service) CheckAuth(req *entity.AuthRequest) (*entity.AuthResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request is nil")
	}

	if req.Login == "" || req.Password == "" || req.IP == "" {
		return nil, fmt.Errorf("login, password and IP must be provided")
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

	// Проверяем ограничения по логину
	loginAllowed, err := s.loginBuckets.Allow(req.Login)
	if err != nil {
		return nil, fmt.Errorf("login check error: %w", err)
	}

	if !loginAllowed {
		return &entity.AuthResponse{OK: false}, nil
	}

	// Проверяем ограничения по паролю
	passwordAllowed, err := s.passwordBuckets.Allow(req.Password)
	if err != nil {
		return nil, fmt.Errorf("password check error: %w", err)
	}

	if !passwordAllowed {
		return &entity.AuthResponse{OK: false}, nil
	}

	// Проверяем ограничения по IP
	ipAllowed, err := s.ipBuckets.Allow(req.IP)
	if err != nil {
		return nil, fmt.Errorf("IP check error: %w", err)
	}

	if !ipAllowed {
		return &entity.AuthResponse{OK: false}, nil
	}

	// Все проверки пройдены
	return &entity.AuthResponse{OK: true}, nil
}

// ResetBucket сбрасывает bucket для указанного логина и IP
func (s *Service) ResetBucket(ctx context.Context, login, ip string) error {
	if login != "" {
		loginKey := fmt.Sprintf("login:%s", login)
		if err := s.loginBuckets.Reset(loginKey); err != nil {
			return fmt.Errorf("login bucket reset error: %w", err)
		}
	}

	if ip != "" {
		ipKey := fmt.Sprintf("ip:%s", ip)
		if err := s.ipBuckets.Reset(ipKey); err != nil {
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
