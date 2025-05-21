package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/vagudza/anti-brute-force/pkg/cli/client"
)

// blacklistCmd represents the blacklist command
var blacklistCmd = &cobra.Command{
	Use:   "blacklist",
	Short: "Manage IP blacklist",
	Long:  `Manage IP blacklist - add, remove and list blacklisted subnets`,
}

// blacklistAddCmd represents the blacklist add command
var blacklistAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add subnet to blacklist",
	Long:  `Add subnet to blacklist in CIDR format (e.g. 192.168.1.0/24)`,
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

		if err := cli.AddToBlacklist(ctx, subnet); err != nil {
			return fmt.Errorf("failed to add subnet to blacklist: %w", err)
		}

		fmt.Printf("Successfully added subnet %q to blacklist\n", subnet)
		return nil
	},
}

// blacklistRemoveCmd represents the blacklist remove command
var blacklistRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove subnet from blacklist",
	Long:  `Remove subnet from blacklist in CIDR format (e.g. 192.168.1.0/24)`,
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

		if err := cli.RemoveFromBlacklist(ctx, subnet); err != nil {
			return fmt.Errorf("failed to remove subnet from blacklist: %w", err)
		}

		fmt.Printf("Successfully removed subnet %q from blacklist\n", subnet)
		return nil
	},
}

// blacklistListCmd represents the blacklist list command
var blacklistListCmd = &cobra.Command{
	Use:   "list",
	Short: "List blacklisted subnets",
	Long:  `List all subnets in blacklist`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cli, err := client.New(host, port)
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		defer cli.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		subnets, err := cli.GetBlacklist(ctx)
		if err != nil {
			return fmt.Errorf("failed to get blacklist: %w", err)
		}

		if len(subnets) == 0 {
			fmt.Println("Blacklist is empty")
			return nil
		}

		fmt.Println("Blacklisted subnets:")
		for _, s := range subnets {
			fmt.Printf("  %s\n", s)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(blacklistCmd)
	blacklistCmd.AddCommand(blacklistAddCmd)
	blacklistCmd.AddCommand(blacklistRemoveCmd)
	blacklistCmd.AddCommand(blacklistListCmd)

	blacklistAddCmd.Flags().StringVar(&subnet, "subnet", "", "Subnet in CIDR format (e.g. 192.168.1.0/24)")
	blacklistRemoveCmd.Flags().StringVar(&subnet, "subnet", "", "Subnet in CIDR format (e.g. 192.168.1.0/24)")
}
