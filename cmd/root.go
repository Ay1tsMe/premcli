package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "premcli",
	Short: "A Premiere League CLI for the terminal.",
	Long: `A Premiere League CLI for the terminal. Displays useful information to track Premiere League games right in the terminal.

Requires an API-FOOTBALL api key found here: https://rapidapi.com/api-sports/api/api-football/`,
}

func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
