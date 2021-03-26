
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Application version",
	Long: `Prints application version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("curfetch v1.0")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
