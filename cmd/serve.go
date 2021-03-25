package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// local command flags
var port int

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve current stored data over http",
	Long: `Serves current database currency data over http, see README for endpoint descriptions`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Serve called on port %v", port)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port for serving API endpoints")
}
