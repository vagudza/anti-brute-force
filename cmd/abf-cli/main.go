package main

import (
	"fmt"
	"os"

	"github.com/vagudza/anti-brute-force/pkg/cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
