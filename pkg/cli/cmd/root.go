package cmd

import (
	"github.com/spf13/cobra"
)

var (
	host string
	port string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "abf-cli",
	Short: "CLI client for Anti-Bruteforce service",
	Long: `CLI client for Anti-Bruteforce service allows you to:
- Reset buckets for specific login/IP
- Manage whitelist/blacklist of IP subnets`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&host, "host", "localhost", "server host")
	rootCmd.PersistentFlags().StringVar(&port, "port", "13013", "server port")
}
