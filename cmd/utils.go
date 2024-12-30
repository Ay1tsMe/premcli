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
	"LEE": "Leeds",
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
	"WBR": "West Brom",
	"WHU": "West Ham",
	"WOL": "Wolves",
}

var teamIDMapping = map[string]int{
	"ARS": 42,   // Arsenal
	"AVL": 66,   // Aston Villa
	"BOU": 35,   // Bournemouth
	"BRE": 55,   // Brentford
	"BHA": 51,   // Brighton
	"BUR": 44,   // Burnley
	"CHE": 49,   // Chelsea
	"CRY": 52,   // Crystal Palace
	"EVE": 45,   // Everton
	"FUL": 36,   // Fulham
	"LEE": 63,   // Leeds
	"LEI": 46,   // Leicester
	"LIV": 40,   // Liverpool
	"LUT": 1359, // Luton
	"MCI": 50,   // Manchester City
	"MUN": 33,   // Manchester United
	"NEW": 34,   // Newcastle
	"NOR": 71,   // Norwich
	"NOT": 65,   // Nottingham Forest
	"SHU": 62,   // Sheffield United
	"SOU": 41,   // Southampton
	"TOT": 47,   // Tottenham
	"WAT": 38,   // Watford
	"WBR": 60,   // West Brom
	"WHU": 48,   // West Ham
	"WOL": 39,   // Wolves
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
