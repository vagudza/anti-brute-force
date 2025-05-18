package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/vagudza/anti-brute-force/api/proto"
	"github.com/vagudza/anti-brute-force/internal/app"
	"github.com/vagudza/anti-brute-force/internal/iplist"
)

func (s *Server) CheckAuth(ctx context.Context, req *pb.CheckAuthRequest) (*pb.CheckAuthResponse, error) {
	authAllowed, err := s.limiterService.CheckAuth(ctx, req.Login, req.Password, req.Ip)
	if err != nil {
		return nil, mapCheckAuthErrors(err)
	}

	if authAllowed {
		return &pb.CheckAuthResponse{Ok: true}, nil
	}

	return &pb.CheckAuthResponse{Ok: false}, nil
}

func (s *Server) ResetBucket(ctx context.Context, req *pb.ResetBucketRequest) (*pb.EmptyResponse, error) {
	if err := s.limiterService.ResetBucket(ctx, req.Login, req.Ip); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.EmptyResponse{}, nil
}

func (s *Server) AddToBlacklist(ctx context.Context, req *pb.IPSubnetRequest) (*pb.EmptyResponse, error) {
	if err := s.limiterService.AddToBlacklist(ctx, req.Subnet); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.EmptyResponse{}, nil
}

func (s *Server) RemoveFromBlacklist(ctx context.Context, req *pb.IPSubnetRequest) (*pb.EmptyResponse, error) {
	if err := s.limiterService.RemoveFromBlacklist(ctx, req.Subnet); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.EmptyResponse{}, nil
}

func (s *Server) AddToWhitelist(ctx context.Context, req *pb.IPSubnetRequest) (*pb.EmptyResponse, error) {
	if err := s.limiterService.AddToWhitelist(ctx, req.Subnet); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.EmptyResponse{}, nil
}

func (s *Server) RemoveFromWhitelist(ctx context.Context, req *pb.IPSubnetRequest) (*pb.EmptyResponse, error) {
	if err := s.limiterService.RemoveFromWhitelist(ctx, req.Subnet); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.EmptyResponse{}, nil
}

func mapCheckAuthErrors(err error) error {
	switch {
	case errors.Is(err, app.ErrEmptyLogin),
		errors.Is(err, app.ErrEmptyPassword),
		errors.Is(err, app.ErrEmptyIP),
		errors.Is(err, app.ErrInvalidIP):
		return status.Error(codes.InvalidArgument, err.Error())

	case errors.Is(err, iplist.ErrInvalidIP),
		errors.Is(err, iplist.ErrInvalidSubnet):
		return status.Error(codes.InvalidArgument, err.Error())

	default:
		return status.Error(codes.Internal, err.Error())
	}
}
