package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var ClearCacheCmd = &cobra.Command{
	Use:   "clear-cache",
	Short: "Clear cache",
	Long:  "Clear the caching proxy's cache",
	Run: func(cmd *cobra.Command, args []string) {
		if err := clearCache(fmt.Sprint(port)); err != nil {
			fmt.Printf("Failed to clear cache: %v\n", err)
			return
		}

		fmt.Printf("Cache cleared successfully\n")
	},
}

func clearCache(port string) error  {
	resp, err := http.Post(fmt.Sprintf("http://localhost:%s/clearCache", port),
		 "",
	nil)
	if err != nil {
		return fmt.Errorf("failed to clear cache: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to clear cache, status code: %d", resp.StatusCode)
	}
	return nil
}

func init() {
	 ClearCacheCmd.Flags().IntVar(&port, "port", 0, "port number")
	_ = ClearCacheCmd.MarkFlagRequired("port")
	rootCmd.AddCommand(ClearCacheCmd)
}
