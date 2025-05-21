package test

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/vagudza/anti-brute-force/api/proto"
	"github.com/vagudza/anti-brute-force/test/suitex"
)

func TestCheckAuth_validations(t *testing.T) {
	ctx, s := suitex.New(t)

	t.Run("empty login", func(t *testing.T) {
		resp, err := s.AntiBrutforceClient.CheckAuth(ctx, &pb.CheckAuthRequest{
			Login:    "",
			Password: "test123",
			Ip:       "192.168.1.1",
		})
		require.Error(t, err)
		require.Nil(t, resp)
		require.Equal(t, codes.InvalidArgument, status.Code(err))
		require.Equal(t, "rpc error: code = InvalidArgument desc = empty login", err.Error())
	})

	t.Run("empty password", func(t *testing.T) {
		resp, err := s.AntiBrutforceClient.CheckAuth(ctx, &pb.CheckAuthRequest{
			Login:    "testuser",
			Password: "",
			Ip:       "192.168.1.1",
		})
		require.Error(t, err)
		require.Nil(t, resp)
		require.Equal(t, codes.InvalidArgument, status.Code(err))
		require.Equal(t, "rpc error: code = InvalidArgument desc = empty password", err.Error())
	})

	t.Run("empty ip", func(t *testing.T) {
		resp, err := s.AntiBrutforceClient.CheckAuth(ctx, &pb.CheckAuthRequest{
			Login:    "testuser",
			Password: "password123",
			Ip:       "",
		})
		require.Error(t, err)
		require.Nil(t, resp)
		require.Equal(t, codes.InvalidArgument, status.Code(err))
		require.Equal(t, "rpc error: code = InvalidArgument desc = empty IP", err.Error())
	})

	t.Run("invalid ip", func(t *testing.T) {
		resp, err := s.AntiBrutforceClient.CheckAuth(ctx, &pb.CheckAuthRequest{
			Login:    "testuser",
			Password: "password123",
			Ip:       "invalid.ip.address",
		})
		require.Error(t, err)
		require.Nil(t, resp)
		require.Equal(t, codes.InvalidArgument, status.Code(err))
		require.Equal(t, "rpc error: code = InvalidArgument desc = invalid IP address", err.Error())
	})

	t.Run("ip with invalid format", func(t *testing.T) {
		resp, err := s.AntiBrutforceClient.CheckAuth(ctx, &pb.CheckAuthRequest{
			Login:    "testuser",
			Password: "password123",
			Ip:       "256.256.256.256",
		})
		require.Error(t, err)
		require.Nil(t, resp)
		require.Equal(t, codes.InvalidArgument, status.Code(err))
		require.Equal(t, "rpc error: code = InvalidArgument desc = invalid IP address", err.Error())
	})
}

func TestCheckAuth(t *testing.T) {
	ctx, s := suitex.New(t)

	t.Run("check auth for specific login", func(t *testing.T) {
		const N = 10 // max attempts per minute for login from config
		login := generateRandomString(t)

		checkAuthWithLogin(ctx, N, t, s, login)
	})

	t.Run("check auth for specific password (reverse brute force)", func(t *testing.T) {
		const M = 10 // max attempts per minute for password from config
		password := generateRandomString(t)

		for i := 0; i < M+1; i++ {
			req := &pb.CheckAuthRequest{
				Login:    generateRandomString(t),
				Password: password,
				Ip:       generateRandomIP(t),
			}

			resp, err := s.AntiBrutforceClient.CheckAuth(ctx, req)
			if i < M {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.True(t, resp.Ok)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.False(t, resp.Ok)
			}
		}
	})

	t.Run("check auth for specific ip", func(t *testing.T) {
		const K = 10 // max attempts per minute for IP from config
		ip := generateRandomIP(t)

		checkAuthWithIP(ctx, K, t, s, ip)
	})

	t.Run("check auth with ip in whitelist", func(t *testing.T) {
		const K = 10 // max attempts per minute for IP from config

		// Generate random IP and create subnet that includes this IP
		ip := generateRandomIP(t)
		// Convert last octet to 0 and add /24 mask to create subnet
		ipParts := strings.Split(ip, ".")
		subnet := fmt.Sprintf("%s.%s.%s.0/24", ipParts[0], ipParts[1], ipParts[2])

		// Add subnet to whitelist
		_, err := s.AntiBrutforceClient.AddToWhitelist(ctx, &pb.IPSubnetRequest{
			Subnet: subnet,
		})
		require.NoError(t, err)

		// Try more requests than allowed by rate limiter
		for i := 0; i < K+10; i++ {
			req := &pb.CheckAuthRequest{
				Login:    generateRandomString(t),
				Password: generateRandomString(t),
				Ip:       ip,
			}

			resp, err := s.AntiBrutforceClient.CheckAuth(ctx, req)
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.True(t, resp.Ok, "Request should be allowed because IP is in whitelist subnet")
		}

		// Cleanup
		_, err = s.AntiBrutforceClient.RemoveFromWhitelist(ctx, &pb.IPSubnetRequest{
			Subnet: subnet,
		})
		require.NoError(t, err)
	})

	t.Run("check auth with ip in blacklist", func(t *testing.T) {
		// Generate random IP and create subnet that includes this IP
		ip := generateRandomIP(t)
		// Convert last octet to 0 and add /24 mask to create subnet
		ipParts := strings.Split(ip, ".")
		subnet := fmt.Sprintf("%s.%s.%s.0/24", ipParts[0], ipParts[1], ipParts[2])

		// Add subnet to blacklist
		_, err := s.AntiBrutforceClient.AddToBlacklist(ctx, &pb.IPSubnetRequest{
			Subnet: subnet,
		})
		require.NoError(t, err)

		req := &pb.CheckAuthRequest{
			Login:    generateRandomString(t),
			Password: generateRandomString(t),
			Ip:       ip,
		}

		resp, err := s.AntiBrutforceClient.CheckAuth(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.False(t, resp.Ok, "Request should be blocked because IP is in blacklist subnet")

		// Cleanup
		_, err = s.AntiBrutforceClient.RemoveFromBlacklist(ctx, &pb.IPSubnetRequest{
			Subnet: subnet,
		})
		require.NoError(t, err)
	})

	t.Run("check auth with specific login and reset bucket", func(t *testing.T) {
		const N = 10 // max attempts per minute for login from config
		login := generateRandomString(t)

		checkAuthWithLogin(ctx, N, t, s, login)

		_, err := s.AntiBrutforceClient.ResetBucket(ctx, &pb.ResetBucketRequest{
			Login: login,
			Ip:    generateRandomIP(t),
		})
		require.NoError(t, err)

		// check auth after reset bucket
		checkAuthWithLogin(ctx, N, t, s, login)
	})

	t.Run("check auth for specific ip and reset bucket", func(t *testing.T) {
		const K = 10 // max attempts per minute for IP from config
		ip := generateRandomIP(t)

		checkAuthWithIP(ctx, K, t, s, ip)

		_, err := s.AntiBrutforceClient.ResetBucket(ctx, &pb.ResetBucketRequest{
			Login: generateRandomString(t),
			Ip:    ip,
		})
		require.NoError(t, err)

		// check auth after reset bucket
		checkAuthWithIP(ctx, K, t, s, ip)
	})
}

func checkAuthWithLogin(
	ctx context.Context,
	N int,
	t *testing.T,
	s *suitex.Suite,
	login string,
) {
	for i := 0; i < N+1; i++ {
		req := &pb.CheckAuthRequest{
			Login:    login,
			Password: generateRandomString(t),
			Ip:       generateRandomIP(t),
		}

		resp, err := s.AntiBrutforceClient.CheckAuth(ctx, req)

		if i < N {
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.True(t, resp.Ok)
		} else {
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.False(t, resp.Ok)
		}
	}
}

func checkAuthWithIP(
	ctx context.Context,
	K int,
	t *testing.T,
	s *suitex.Suite,
	ip string,
) {
	for i := 0; i < K+1; i++ {
		req := &pb.CheckAuthRequest{
			Login:    generateRandomString(t),
			Password: generateRandomString(t),
			Ip:       ip,
		}

		resp, err := s.AntiBrutforceClient.CheckAuth(ctx, req)
		if i < K {
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.True(t, resp.Ok)
		} else {
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.False(t, resp.Ok)
		}
	}
}

func generateRandomString(t *testing.T) string {
	b := make([]byte, 20)
	_, err := rand.Read(b)
	require.NoError(t, err)
	return base64.StdEncoding.EncodeToString(b)
}

func generateRandomIP(t *testing.T) string {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	require.NoError(t, err)
	return fmt.Sprintf("%d.%d.%d.%d", b[0], b[1], b[2], b[3])
}
