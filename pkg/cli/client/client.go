package client

import (
	"context"
	"fmt"

	pb "github.com/vagudza/anti-brute-force/api/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.AntiBruteforceClient
}

func New(host, port string) (*Client, error) {
	addr := fmt.Sprintf("%s:%s", host, port)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	return &Client{
		conn:   conn,
		client: pb.NewAntiBruteforceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

// ResetBucket resets bucket for login and IP
func (c *Client) ResetBucket(ctx context.Context, login, ip string) error {
	_, err := c.client.ResetBucket(ctx, &pb.ResetBucketRequest{
		Login: login,
		Ip:    ip,
	})
	return err
}

// AddToWhitelist adds subnet to whitelist
func (c *Client) AddToWhitelist(ctx context.Context, subnet string) error {
	_, err := c.client.AddToWhitelist(ctx, &pb.IPSubnetRequest{
		Subnet: subnet,
	})
	return err
}

// RemoveFromWhitelist removes subnet from whitelist
func (c *Client) RemoveFromWhitelist(ctx context.Context, subnet string) error {
	_, err := c.client.RemoveFromWhitelist(ctx, &pb.IPSubnetRequest{
		Subnet: subnet,
	})
	return err
}

// GetWhitelist returns all subnets from whitelist
func (c *Client) GetWhitelist(ctx context.Context) ([]string, error) {
	resp, err := c.client.GetWhitelist(ctx, &pb.EmptyRequest{})
	if err != nil {
		return nil, err
	}
	return resp.Subnets, nil
}

// AddToBlacklist adds subnet to blacklist
func (c *Client) AddToBlacklist(ctx context.Context, subnet string) error {
	_, err := c.client.AddToBlacklist(ctx, &pb.IPSubnetRequest{
		Subnet: subnet,
	})
	return err
}

// RemoveFromBlacklist removes subnet from blacklist
func (c *Client) RemoveFromBlacklist(ctx context.Context, subnet string) error {
	_, err := c.client.RemoveFromBlacklist(ctx, &pb.IPSubnetRequest{
		Subnet: subnet,
	})
	return err
}

// GetBlacklist returns all subnets from blacklist
func (c *Client) GetBlacklist(ctx context.Context) ([]string, error) {
	resp, err := c.client.GetBlacklist(ctx, &pb.EmptyRequest{})
	if err != nil {
		return nil, err
	}
	return resp.Subnets, nil
}
