package test

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	pb "github.com/vagudza/anti-brute-force/api/proto"
	"github.com/vagudza/anti-brute-force/test/suitex"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCheckAuth_validations(t *testing.T) {
	ctx, s := suitex.New(t)
	clearLists(ctx, t, s)

	t.Run("empty login", func(t *testing.T) {
		resp, err := s.AntiBruteforceClient.CheckAuth(ctx, &pb.CheckAuthRequest{
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
		resp, err := s.AntiBruteforceClient.CheckAuth(ctx, &pb.CheckAuthRequest{
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
		resp, err := s.AntiBruteforceClient.CheckAuth(ctx, &pb.CheckAuthRequest{
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
		resp, err := s.AntiBruteforceClient.CheckAuth(ctx, &pb.CheckAuthRequest{
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
		resp, err := s.AntiBruteforceClient.CheckAuth(ctx, &pb.CheckAuthRequest{
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
	clearLists(ctx, t, s)

	// set values from app config
	N := s.Cfg.Limiters.Login.MaxAttemptsPerMinute
	M := s.Cfg.Limiters.Password.MaxAttemptsPerMinute
	K := s.Cfg.Limiters.IP.MaxAttemptsPerMinute

	t.Run("check auth for specific login", func(t *testing.T) {
		login := generateRandomString(t)
		checkAuthWithLogin(ctx, t, N, s, login)
	})

	t.Run("check auth for specific password (reverse brute force)", func(t *testing.T) {
		password := generateRandomString(t)
		checkAuthWithPassword(ctx, t, M, s, password)
	})

	t.Run("check auth for specific ip", func(t *testing.T) {
		ip := generateRandomIP(t)
		checkAuthWithIP(ctx, t, K, s, ip)
	})

	t.Run("check auth with ip in whitelist", func(t *testing.T) {
		// Generate random IP and create subnet that includes this IP
		ip := generateRandomIP(t)

		// Convert last octet to 0 and add /24 mask to create subnet
		ipParts := strings.Split(ip, ".")
		subnet := fmt.Sprintf("%s.%s.%s.0/24", ipParts[0], ipParts[1], ipParts[2])

		// Add subnet to whitelist
		_, err := s.AntiBruteforceClient.AddToWhitelist(ctx, &pb.IPSubnetRequest{
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

			resp, errResp := s.AntiBruteforceClient.CheckAuth(ctx, req)
			require.NoError(t, errResp)
			require.NotNil(t, resp)
			require.True(t, resp.Ok, "Request should be allowed because IP is in whitelist subnet")
		}

		// Cleanup
		_, err = s.AntiBruteforceClient.RemoveFromWhitelist(ctx, &pb.IPSubnetRequest{
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
		_, err := s.AntiBruteforceClient.AddToBlacklist(ctx, &pb.IPSubnetRequest{
			Subnet: subnet,
		})
		require.NoError(t, err)

		req := &pb.CheckAuthRequest{
			Login:    generateRandomString(t),
			Password: generateRandomString(t),
			Ip:       ip,
		}

		resp, err := s.AntiBruteforceClient.CheckAuth(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.False(t, resp.Ok, "Request should be blocked because IP is in blacklist subnet")

		// Cleanup
		_, err = s.AntiBruteforceClient.RemoveFromBlacklist(ctx, &pb.IPSubnetRequest{
			Subnet: subnet,
		})
		require.NoError(t, err)
	})

	t.Run("check auth with specific login and reset bucket", func(t *testing.T) {
		login := generateRandomString(t)
		checkAuthWithLogin(ctx, t, N, s, login)

		_, err := s.AntiBruteforceClient.ResetBucket(ctx, &pb.ResetBucketRequest{
			Login: login,
			Ip:    generateRandomIP(t),
		})
		require.NoError(t, err)

		// check auth after reset bucket
		checkAuthWithLogin(ctx, t, N, s, login)
	})

	t.Run("check auth for specific ip and reset bucket", func(t *testing.T) {
		ip := generateRandomIP(t)
		checkAuthWithIP(ctx, t, K, s, ip)

		_, err := s.AntiBruteforceClient.ResetBucket(ctx, &pb.ResetBucketRequest{
			Login: generateRandomString(t),
			Ip:    ip,
		})
		require.NoError(t, err)

		// check auth after reset bucket
		checkAuthWithIP(ctx, t, K, s, ip)
	})
}

func clearLists(
	ctx context.Context,
	t *testing.T,
	s *suitex.Suite,
) {
	t.Helper()
	t.Log("clearing whitelist and blacklist before test")

	_, err := s.AntiBruteforceClient.ClearBlacklist(ctx, &pb.EmptyRequest{})
	require.NoError(t, err)
	_, err = s.AntiBruteforceClient.ClearWhitelist(ctx, &pb.EmptyRequest{})
	require.NoError(t, err)
}

func checkAuthWithLogin(
	ctx context.Context,
	t *testing.T,
	n int,
	s *suitex.Suite,
	login string,
) {
	t.Helper()

	for i := 0; i < n+1; i++ {
		req := &pb.CheckAuthRequest{
			Login:    login,
			Password: generateRandomString(t),
			Ip:       generateRandomIP(t),
		}

		resp, err := s.AntiBruteforceClient.CheckAuth(ctx, req)

		if i < n {
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.True(t, resp.Ok) // Expecting success for first n attempts
		} else {
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.False(t, resp.Ok) // Expecting failure for the (n+1)th attempt
		}
	}
}

func checkAuthWithPassword(
	ctx context.Context,
	t *testing.T,
	m int,
	s *suitex.Suite,
	password string,
) {
	t.Helper()
	for i := 0; i < m+1; i++ {
		req := &pb.CheckAuthRequest{
			Login:    generateRandomString(t),
			Password: password,
			Ip:       generateRandomIP(t),
		}

		resp, err := s.AntiBruteforceClient.CheckAuth(ctx, req)
		if i < m {
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
	t *testing.T,
	k int,
	s *suitex.Suite,
	ip string,
) {
	t.Helper()

	for i := 0; i < k+1; i++ {
		req := &pb.CheckAuthRequest{
			Login:    generateRandomString(t),
			Password: generateRandomString(t),
			Ip:       ip,
		}

		resp, err := s.AntiBruteforceClient.CheckAuth(ctx, req)
		if i < k {
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
	t.Helper()

	b := make([]byte, 20)
	_, err := rand.Read(b)
	require.NoError(t, err)
	return base64.StdEncoding.EncodeToString(b)
}

func generateRandomIP(t *testing.T) string {
	t.Helper()

	b := make([]byte, 4)
	_, err := rand.Read(b)
	require.NoError(t, err)
	return fmt.Sprintf("%d.%d.%d.%d", b[0], b[1], b[2], b[3])
}

func generateRandomSubnet(t *testing.T) string {
	t.Helper()

	// Generate random IP
	b := make([]byte, 4)
	_, err := rand.Read(b)
	require.NoError(t, err)

	// Generate random CIDR mask between 8 and 30 using crypto/rand
	maskBytes := make([]byte, 8)
	_, err = rand.Read(maskBytes)
	require.NoError(t, err)

	// Convert to int64 and get value in range 8-30
	mask := int64(maskBytes[0])%23 + 8 // 8-30

	bitsToZero := 32 - mask
	bytesToZero := bitsToZero / 8
	remainingBits := bitsToZero % 8

	for i := 4 - bytesToZero; i < 4; i++ {
		b[i] = 0
	}

	if remainingBits > 0 && 4-bytesToZero-1 >= 0 {
		idx := 4 - bytesToZero - 1
		mask := byte(0xFF << remainingBits)
		b[idx] &= mask
	}

	return fmt.Sprintf("%d.%d.%d.%d/%d", b[0], b[1], b[2], b[3], mask)
}
