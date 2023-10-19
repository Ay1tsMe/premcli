/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	apiKey   string
	timezone string
	favTeam  string
)

type ApiResponse struct {
	Response []Match `json:"response"`
}

type Match struct {
	Teams struct {
		Home struct {
			Name string
		}
		Away struct {
			Name string
		}
	}
}

func GetConfig() error {
	configFile, err := os.Open(configPath)
	if err != nil {
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

func getCurrentRound() string {
	return "TODO"
}

func buildURL() string {
	baseURL := "https://api-football-v1.p.rapidapi.com/v3/fixtures?league=39&season=2023&round=Regular%20Season%20-%209&timezone="
	return baseURL + url.QueryEscape(timezone)
}

func fetchAndParse() ([]Match, error) {
	url := buildURL()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %v", err)
	}

	req.Header.Add("X-RapidAPI-Key", apiKey)
	req.Header.Add("X-RapidAPI-Host", "api-football-v1.p.rapidapi.com")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error executing request: %v", err)
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response: %v", err)
	}

	var responseData ApiResponse
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return nil, fmt.Errorf("Error parsing JSON response: %v", err)
	}

	return responseData.Response, nil
}

// fixturesCmd represents the fixtures command
var fixturesCmd = &cobra.Command{
	Use:   "fixtures",
	Short: "Prints fixtures for current round",
	Long:  "Prints fixtures for current round and highlights in bold the fixture of your favourite team.",
	Run: func(cmd *cobra.Command, args []string) {
		err := GetConfig()
		if err != nil {
			fmt.Println("Error loading config:", err)
			return
		}

		matches, err := fetchAndParse()
		if err != nil {
			fmt.Println("Error fetching and parsing:", err)
			return
		}

		for _, match := range matches {
			fmt.Printf("[H] %s vs. %s [A]\n", match.Teams.Home.Name, match.Teams.Away.Name)
		}
	},
}

func init() {
	rootCmd.AddCommand(fixturesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fixturesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fixturesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
