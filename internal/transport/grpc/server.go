package grpc

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/vagudza/anti-brute-force/api/proto"
	"github.com/vagudza/anti-brute-force/internal/app"
	"github.com/vagudza/anti-brute-force/internal/config"
)

// Server представляет собой обертку над gRPC сервером
type Server struct {
	server         *grpc.Server
	listener       net.Listener
	cfg            *config.GrpcConfig
	limiterService app.LimiterService

	pb.UnimplementedAntiBruteforceServer
}

// NewServer создает новый gRPC сервер
func NewServer(limiterService app.LimiterService, cfg *config.GrpcConfig) *Server {
	s := &Server{
		server:         grpc.NewServer(),
		cfg:            cfg,
		limiterService: limiterService,
	}
	pb.RegisterAntiBruteforceServer(s.server, s)
	reflection.Register(s.server)
	return s
}

// Start запускает gRPC сервер на указанном адресе
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.cfg.Port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s.listener = listener

	// Запускаем gRPC сервер
	if err := s.server.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}

// Stop останавливает gRPC сервер с грациозным завершением
func (s *Server) Stop(ctx context.Context) error {
	// Создаем канал для отслеживания завершения остановки
	stopped := make(chan struct{})

	go func() {
		s.server.GracefulStop()
		close(stopped)
	}()

	// Ожидаем либо завершения остановки, либо истечения таймаута
	select {
	case <-ctx.Done():
		// Если контекст истек, принудительно останавливаем сервер
		s.server.Stop()
		return fmt.Errorf("server shutdown timed out")
	case <-stopped:
		// Сервер успешно остановлен
		return nil
	}
}
