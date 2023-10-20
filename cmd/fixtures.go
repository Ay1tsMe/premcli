/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	apiKey     string
	timezone   string
	favTeam    string
	roundValue string
)

type ApiResponse struct {
	Response []Match `json:"response"`
}

type Match struct {
	Fixture struct {
		Date string `json:"date"`
	} `json:"fixture"`
	Teams struct {
		Home struct {
			Name string
		}
		Away struct {
			Name string
		}
	}
}

type CurrentRound struct {
	Response []string `json:"response"`
}

func getCurrentRound() error {
	url := buildRoundURL()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("Error creating request: %v", err)
	}

	req.Header.Add("X-RapidAPI-Key", apiKey)
	req.Header.Add("X-RapidAPI-Host", "api-football-v1.p.rapidapi.com")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Error executing request: %v", err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Error reasing response: %v", err)
	}

	var responseData CurrentRound
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return fmt.Errorf("Error parsing JSON response: %v", err)
	}

	if len(responseData.Response) > 0 {
		roundValue = responseData.Response[0]
		return nil
	}

	return fmt.Errorf("No round information found in the API response")
}

func buildURL() string {
	baseURL := "https://api-football-v1.p.rapidapi.com/v3/fixtures?league=39"

	season := "&season=" + getSeasonYear()

	round := "&round=" + url.QueryEscape(roundValue)
	tz := "&timezone=" + url.QueryEscape(timezone)

	return baseURL + season + round + tz
}

func buildRoundURL() string {
	baseURL := "https://api-football-v1.p.rapidapi.com/v3/fixtures/rounds?league=39&current=true"

	season := "&season=" + getSeasonYear()

	return baseURL + season
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
	body, err := io.ReadAll(res.Body)
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

func FormatTime(isoTime string) (string, error) {
	parsedTime, err := time.Parse(time.RFC3339, isoTime)
	if err != nil {
		return "", fmt.Errorf("Error parsing time: %v", err)
	}

	return parsedTime.Format("02 Jan 2006, 03:04 PM"), nil
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

		err = getCurrentRound()
		if err != nil {
			fmt.Println("Error getting current round:", err)
			return
		}

		matches, err := fetchAndParse()
		if err != nil {
			fmt.Println("Error fetching and parsing:", err)
			return
		}

		color.Set(color.Underline)
		fmt.Println(roundValue)
		color.Unset()

		for _, match := range matches {
			homeTeam := match.Teams.Home.Name
			awayTeam := match.Teams.Away.Name
			date := match.Fixture.Date

			userFriendlyTime, err := FormatTime(date)
			if err != nil {
				fmt.Println(err)
				return
			}

			matchDisplay := fmt.Sprintf("[H] %s vs. %s [A]\nDate: %s", homeTeam, awayTeam, userFriendlyTime)

			if isFavTeam(homeTeam, awayTeam, favTeam) {
				color.Set(color.Bold)
				fmt.Println(matchDisplay)
				color.Unset()
				continue
			}

			fmt.Println(matchDisplay)
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
