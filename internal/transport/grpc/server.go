package grpc

import (
	"context"
	"fmt"
	"net"

	pb "github.com/vagudza/anti-brute-force/api/proto"
	"github.com/vagudza/anti-brute-force/internal/app"
	"github.com/vagudza/anti-brute-force/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", ":"+s.cfg.Port)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s.listener = listener
	if err = s.server.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
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
