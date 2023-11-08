/*
Contains utlities that can be used across all .go files
*/
package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

var teamMapping = map[string]string{
	"ARS": "Arsenal",
	"AVL": "Aston Villa",
	"BOU": "Bournemouth",
	"BRE": "Brentford",
	"BHA": "Brighton",
	"BUR": "Burnley",
	"CHE": "Chelsea",
	"CRY": "Crystal Palace",
	"EVE": "Everton",
	"FUL": "Fulham",
	"LEI": "Leicester City",
	"LIV": "Liverpool",
	"LUT": "Luton",
	"MCI": "Manchester City",
	"MUN": "Manchester United",
	"NEW": "Newcastle",
	"NOR": "Norwich City",
	"NOT": "Nottingham Forest",
	"SHU": "Sheffield United",
	"SOU": "Southampton",
	"TOT": "Tottenham Hotspur",
	"WAT": "Watford",
	"WHU": "West Ham",
	"WOL": "Wolves",
}

var (
	apiKey   string
	timezone string
	favTeam  string
)

// Retrieves the config information
func GetConfig() error {
	configFile, err := os.Open(configPath)
	if err != nil {
		var pathErr *os.PathError
		if errors.As(err, &pathErr) {
			return fmt.Errorf("Please run 'premcli config' to set up a configuration.")
		}
		return fmt.Errorf("Failed to open config file: %v", err)
	}
	defer configFile.Close()

	scanner := bufio.NewScanner(configFile)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)

		if len(parts) != 2 {
			return fmt.Errorf("Invalid config line: %s", line)
		}

		key := parts[0]
		value := parts[1]

		switch key {
		case "API_KEY":
			apiKey = value
		case "TIMEZONE":
			timezone = value
		case "FAVTEAM":
			favTeam = value
		default:
			return fmt.Errorf("Unknown config key: %s", key)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Error reading config file: %v", err)
	}

	return nil
}

// Checks if favourite team is present within string
func isFavTeam(matchString, favTeam string) bool {
	teamName, exists := teamMapping[strings.ToUpper(favTeam)]
	if !exists {
		return false
	}

	return strings.Contains(strings.ToUpper(matchString), strings.ToUpper(teamName))
}

// Gets the current season year
func getSeasonYear() string {
	currentYear := time.Now().Year()
	currentMonth := time.Now().Month()

	if currentMonth < time.July {
		currentYear--
	}

	return fmt.Sprintf("%d", currentYear)
}
