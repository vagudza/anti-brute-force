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

type Server struct {
	server         *grpc.Server
	listener       net.Listener
	cfg            *config.GrpcConfig
	limiterService app.LimiterService

	pb.UnimplementedAntiBruteforceServer
}

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
	listener, err := net.Listen("tcp", ":"+s.cfg.Port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s.listener = listener
	if err = s.server.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	stopped := make(chan struct{})

	go func() {
		s.server.GracefulStop()
		close(stopped)
	}()

	select {
	case <-ctx.Done():
		s.server.Stop()
		return fmt.Errorf("server shutdown timed out")
	case <-stopped:
		return nil
	}
}
