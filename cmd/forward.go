package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

var (
    port   int
    origin string
)

var forwardCmd = &cobra.Command{
    Use:     "forward",
    Short:   "Forward request",
    Long:    "Forward request from one URL to another",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Printf("forwarding request from port %d to %s \n", port, origin)
        Server(fmt.Sprint(port), origin)
    },
}

func init() {
    forwardCmd.Flags().IntVar(&port, "port", 0, "port number")
    forwardCmd.Flags().StringVar(&origin, "origin", "", "origin url")
    _ = forwardCmd.MarkFlagRequired("port")
    _ = forwardCmd.MarkFlagRequired("origin")
    rootCmd.AddCommand(forwardCmd)
}