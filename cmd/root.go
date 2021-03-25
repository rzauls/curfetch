package cmd

import (
	"github.com/spf13/cobra"
)

// rootCmd base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "curfetch",
	Short: "Currency RSS feed reader and http api",
	Long: "This application fetches currency data from a specific RSS feed, uploads them to a database and can serve them as http endpoints",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
