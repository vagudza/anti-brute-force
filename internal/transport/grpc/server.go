package grpc

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server представляет собой обертку над gRPC сервером
type Server struct {
	server   *grpc.Server
	listener net.Listener
	port     string
}

// Service интерфейс, который должны реализовывать все gRPC сервисы
type Service interface {
	Register(*grpc.Server)
}

// NewServer создает новый gRPC сервер
func NewServer(service Service, port string) *Server {
	server := grpc.NewServer()

	// Регистрируем сервис в gRPC сервере
	service.Register(server)

	// Включаем reflection API для упрощения отладки с помощью таких инструментов, как grpcurl
	reflection.Register(server)

	return &Server{
		server: server,
		port:   port,
	}
}

// Start запускает gRPC сервер на указанном адресе
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.port)
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
