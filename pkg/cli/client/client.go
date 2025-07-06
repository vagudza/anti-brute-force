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

func (c *Client) ResetBucket(ctx context.Context, login, ip string) error {
	_, err := c.client.ResetBucket(ctx, &pb.ResetBucketRequest{
		Login: login,
		Ip:    ip,
	})
	return err
}

func (c *Client) AddToWhitelist(ctx context.Context, subnet string) error {
	_, err := c.client.AddToWhitelist(ctx, &pb.IPSubnetRequest{
		Subnet: subnet,
	})
	return err
}

func (c *Client) RemoveFromWhitelist(ctx context.Context, subnet string) error {
	_, err := c.client.RemoveFromWhitelist(ctx, &pb.IPSubnetRequest{
		Subnet: subnet,
	})
	return err
}

func (c *Client) GetWhitelist(ctx context.Context) ([]string, error) {
	resp, err := c.client.GetWhitelist(ctx, &pb.EmptyRequest{})
	if err != nil {
		return nil, err
	}
	return resp.Subnets, nil
}

func (c *Client) AddToBlacklist(ctx context.Context, subnet string) error {
	_, err := c.client.AddToBlacklist(ctx, &pb.IPSubnetRequest{
		Subnet: subnet,
	})
	return err
}

func (c *Client) RemoveFromBlacklist(ctx context.Context, subnet string) error {
	_, err := c.client.RemoveFromBlacklist(ctx, &pb.IPSubnetRequest{
		Subnet: subnet,
	})
	return err
}

func (c *Client) GetBlacklist(ctx context.Context) ([]string, error) {
	resp, err := c.client.GetBlacklist(ctx, &pb.EmptyRequest{})
	if err != nil {
		return nil, err
	}
	return resp.Subnets, nil
}
