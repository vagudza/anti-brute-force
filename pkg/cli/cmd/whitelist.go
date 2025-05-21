package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/vagudza/anti-brute-force/pkg/cli/client"
)

var subnet string

// whitelistCmd represents the whitelist command
var whitelistCmd = &cobra.Command{
	Use:   "whitelist",
	Short: "Manage IP whitelist",
	Long:  `Manage IP whitelist - add, remove and list whitelisted subnets`,
}

// whitelistAddCmd represents the whitelist add command
var whitelistAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add subnet to whitelist",
	Long:  `Add subnet to whitelist in CIDR format (e.g. 192.168.1.0/24)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if subnet == "" {
			return fmt.Errorf("subnet must be specified")
		}

		cli, err := client.New(host, port)
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		defer cli.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := cli.AddToWhitelist(ctx, subnet); err != nil {
			return fmt.Errorf("failed to add subnet to whitelist: %w", err)
		}

		fmt.Printf("Successfully added subnet %q to whitelist\n", subnet)
		return nil
	},
}

// whitelistRemoveCmd represents the whitelist remove command
var whitelistRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove subnet from whitelist",
	Long:  `Remove subnet from whitelist in CIDR format (e.g. 192.168.1.0/24)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if subnet == "" {
			return fmt.Errorf("subnet must be specified")
		}

		cli, err := client.New(host, port)
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		defer cli.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := cli.RemoveFromWhitelist(ctx, subnet); err != nil {
			return fmt.Errorf("failed to remove subnet from whitelist: %w", err)
		}

		fmt.Printf("Successfully removed subnet %q from whitelist\n", subnet)
		return nil
	},
}

// whitelistListCmd represents the whitelist list command
var whitelistListCmd = &cobra.Command{
	Use:   "list",
	Short: "List whitelisted subnets",
	Long:  `List all subnets in whitelist`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cli, err := client.New(host, port)
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		defer cli.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		subnets, err := cli.GetWhitelist(ctx)
		if err != nil {
			return fmt.Errorf("failed to get whitelist: %w", err)
		}

		if len(subnets) == 0 {
			fmt.Println("Whitelist is empty")
			return nil
		}

		fmt.Println("Whitelisted subnets:")
		for _, s := range subnets {
			fmt.Printf("  %s\n", s)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(whitelistCmd)
	whitelistCmd.AddCommand(whitelistAddCmd)
	whitelistCmd.AddCommand(whitelistRemoveCmd)
	whitelistCmd.AddCommand(whitelistListCmd)

	whitelistAddCmd.Flags().StringVar(&subnet, "subnet", "", "Subnet in CIDR format (e.g. 192.168.1.0/24)")
	whitelistRemoveCmd.Flags().StringVar(&subnet, "subnet", "", "Subnet in CIDR format (e.g. 192.168.1.0/24)")
}
