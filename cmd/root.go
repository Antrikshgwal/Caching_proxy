package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "caching-proxy",
	Short: "A cli tool that starts a proxy cache server.",
	Long:  "A cli tool that starts a proxy cache server. It will forward requests to the actual server and cache the responses. If the same request is made again, it will return the cached response instead of forwarding the request to the server.",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Oops. An error while executing caching-proxy '%s'\n", err)
		os.Exit(1)
	}
}
