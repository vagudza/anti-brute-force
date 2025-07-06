package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.PersistentFlags().StringVar(&host, "host", "localhost", "server host")
	rootCmd.PersistentFlags().StringVar(&port, "port", "13013", "server port")
}

var (
	host string
	port string
)

var rootCmd = &cobra.Command{
	Use:   "abf-cli",
	Short: "CLI client for Anti-Bruteforce service",
	Long: `CLI client for Anti-Bruteforce service allows you to:
- Reset buckets for specific login/IP
- Manage whitelist/blacklist of IP subnets`,
	SilenceUsage:  true, // skip usage when error
	SilenceErrors: true, // show only our errors
}

func Execute() error {
	return rootCmd.Execute()
}
