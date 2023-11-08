/*
Sets up premcli.conf file for the user
Collects information such as:
- API Key
- TimeZone
- Favourite Team
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var configPath = filepath.Join(os.Getenv("HOME"), ".config", "premcli", "premcli.conf")
var overwrite bool

// Checks if premcli.conf exists
func ConfigExists() bool {
	if _, err := os.Stat(configPath); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		return false
	}
}

// Creates the premcli directory
func CreateDir() bool {
	dir := filepath.Dir(configPath)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Printf("Failed to create directory %s\n", err)
			return false
		}
	}
	return true
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure settings for premcli",
	Run: func(cmd *cobra.Command, args []string) {
		// Deletes old premcli.conf if overwrite flag is called
		if overwrite {
			err := os.Remove(configPath)
			if err != nil {
				fmt.Printf("Failed to remove existing config file: %s\n", err)
				return
			}
			fmt.Println("Existing config file removed.")
		}

		// Return if premcli.conf exists
		if ConfigExists() {
			fmt.Println("Config file already exists at", configPath)
			fmt.Println("Use the -overwrite flag to replace it.")
			return
		}

		// Get config variables
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Welcome to premcli. To use premcli, you must acquire an api key from API-FOOTBALL found here: https://rapidapi.com/api-sports/api/api-football/\n")
		fmt.Print("Enter your API key: ")
		apiKey, _ := reader.ReadString('\n')
		apiKey = strings.TrimSpace(apiKey)

		fmt.Print("TIMEZONE EXAMPLE: Europe/Berlin, Australia/Sydney\n")
		fmt.Print("Use 'premcli timezones' to list all available timezones.\n")
		fmt.Print("Enter your Timezone: ")
		timezone, _ := reader.ReadString('\n')
		timezone = strings.TrimSpace(timezone)

		fmt.Print("TEAM EXAMPLE: WOL, MCI, CHE\n")
		fmt.Print("Enter your Premier League team: ")
		favTeam, _ := reader.ReadString('\n')
		favTeam = strings.TrimSpace(favTeam)

		content := fmt.Sprintf("API_KEY=%s\nTIMEZONE=%s\nFAVTEAM=%s\n", apiKey, timezone, favTeam)

		if !CreateDir() {
			return
		}

		// Writes variables to premcli.conf
		err := os.WriteFile(configPath, []byte(content), 0644)
		if err != nil {
			fmt.Printf("Failed to write to config file: %s\n", err)
			return
		}

		fmt.Println("Configuration saved!")
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().BoolVarP(&overwrite, "overwrite", "o", false, "Overwrite the existing config")
}
