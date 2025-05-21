package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "1.0.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Print the version number of abf-cli`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("abf-cli version %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
