package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/vagudza/anti-brute-force/pkg/cli/client"
)

var (
	login string
	ip    string
)

// bucketCmd represents the bucket command
var bucketCmd = &cobra.Command{
	Use:   "bucket",
	Short: "Manage rate limit buckets",
	Long:  `Manage rate limit buckets for specific login/IP combinations`,
}

// resetCmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset bucket for login/IP",
	Long:  `Reset bucket for specific login/IP combination`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if login == "" && ip == "" {
			return fmt.Errorf("either login or IP must be specified")
		}

		cli, err := client.New(host, port)
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		defer cli.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := cli.ResetBucket(ctx, login, ip); err != nil {
			return fmt.Errorf("failed to reset bucket: %w", err)
		}

		fmt.Printf("Successfully reset bucket for login=%q ip=%q\n", login, ip)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(bucketCmd)
	bucketCmd.AddCommand(resetCmd)

	resetCmd.Flags().StringVar(&login, "login", "", "Login to reset bucket for")
	resetCmd.Flags().StringVar(&ip, "ip", "", "IP to reset bucket for")
}
