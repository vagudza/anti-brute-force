package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/vagudza/anti-brute-force/api/proto"
	"github.com/vagudza/anti-brute-force/internal/app"
)

// type AntiBruteforceServer interface {
// 	// CheckAuth verifies an authentication attempt
// 	CheckAuth(context.Context, *CheckAuthRequest) (*CheckAuthResponse, error)
// 	// ResetBucket resets the bucket by login and IP
// 	ResetBucket(context.Context, *ResetBucketRequest) (*EmptyResponse, error)
// 	// AddToBlacklist adds a subnet to the blacklist
// 	AddToBlacklist(context.Context, *IPSubnetRequest) (*EmptyResponse, error)
// 	// RemoveFromBlacklist removes a subnet from the blacklist
// 	RemoveFromBlacklist(context.Context, *IPSubnetRequest) (*EmptyResponse, error)
// 	// AddToWhitelist adds a subnet to the whitelist
// 	AddToWhitelist(context.Context, *IPSubnetRequest) (*EmptyResponse, error)
// 	// RemoveFromWhitelist removes a subnet from the whitelist
// 	RemoveFromWhitelist(context.Context, *IPSubnetRequest) (*EmptyResponse, error)
// 	// GetBlacklist retrieves all subnets from the blacklist
// 	GetBlacklist(context.Context, *EmptyRequest) (*IPSubnetListResponse, error)
// 	// GetWhitelist retrieves all subnets from the whitelist
// 	GetWhitelist(context.Context, *EmptyRequest) (*IPSubnetListResponse, error)
// 	mustEmbedUnimplementedAntiBruteforceServer()
// }

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

// ResetBucket resets the bucket by login and IP
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
	case errors.Is(err, app.ErrEmptyLogin), errors.Is(err, app.ErrEmptyPassword), errors.Is(err, app.ErrEmptyIP):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
