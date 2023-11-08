/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
)

type ApiResponseEvents struct {
	Response []Events `json:"response"`
}

type Events struct {
	Time struct {
		Elapsed int
		Extra   int
	}
	Team struct {
		Name string
	}
	Player struct {
		Name string
	}
	Assist struct {
		Name string
	}
	Type     string
	Detail   string
	Comments string
}

type ApiResponseFixtureByID struct {
	Response []Match `json:"response"`
}

func buildFixtureByIDURL(fixtureID int) string {
	baseURL := "https://api-football-v1.p.rapidapi.com/v3/fixtures?"

	fixture := "id=" + strconv.Itoa(fixtureID)

	return baseURL + fixture
}

func buildEventsURL(fixtureID int) string {
	baseURL := "https://api-football-v1.p.rapidapi.com/v3/fixtures/events?"

	fixture := "fixture=" + strconv.Itoa(fixtureID)

	return baseURL + fixture
}

func getEvents(fixtureID int) ([]Events, error) {
	url := buildEventsURL(fixtureID)

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

	var responseData ApiResponseEvents
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return nil, fmt.Errorf("Error parsing JSON response: %v", err)
	}

	return responseData.Response, nil
}

func getFixtureByID(fixtureID int) ([]Match, error) {
	url := buildFixtureByIDURL(fixtureID)

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

	var responseData ApiResponseFixtureByID
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return nil, fmt.Errorf("Error parsing JSON response: %v", err)
	}

	return responseData.Response, nil
}

// liveCmd represents the live command
var liveCmd = &cobra.Command{
	Use:   "live",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := GetConfig()
		if err != nil {
			fmt.Println("Error loading config:", err)
			return
		}

		fixtureID, err := strconv.Atoi(args[0])
		if err != nil {
			return
		}

		events, err := getEvents(fixtureID)
		if err != nil {
			fmt.Println("Error fetching and parsing:", err)
			return
		}

		match, err := getFixtureByID(fixtureID)
		if err != nil {
			fmt.Println("Error fetching and parsing:", err)
			return
		}

		homeTeam := match[0].Teams.Home.Name
		homeScore := match[0].Goals.Home
		awayTeam := match[0].Teams.Away.Name
		awayScore := match[0].Goals.Away
		date := match[0].Fixture.Date
		timeElapsed := match[0].Fixture.Status.Elapsed

		// Reformat time so its readable
		userFriendlyTime, err := FormatTime(date)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Score Padding
		const nameScoreWidth = 26
		const vsWidth = 25

		homePadding := nameScoreWidth - len("[H]") - len(homeTeam) - len(fmt.Sprint(homeScore))
		awayPadding := nameScoreWidth - len("[A]") - len(awayTeam) - len(fmt.Sprint(awayScore))

		// Fixture Info
		matchDisplay := ""
		matchDisplay = fmt.Sprintf("Date: %s\n[H] %s%*s%d\n[A] %s%*s%d\nTime Elapsed: %d\nEvents:\n", userFriendlyTime, homeTeam, homePadding, "", homeScore, awayTeam, awayPadding, "", awayScore, timeElapsed)

		// Events Info
		var eventsArr []string

		for _, event := range events {
			eventTime := event.Time.Elapsed
			teamName := event.Team.Name
			playerName := event.Player.Name
			assistName := event.Assist.Name
			eventType := event.Type
			eventDetail := event.Detail
			eventComment := event.Comments

			eventSummary := ""
			if eventType == "Card" {
				eventSummary = fmt.Sprintf("%d' %s\n%s\n%s\n%s\n", eventTime, eventDetail, teamName, playerName, eventComment)
			} else if eventType == "subst" {
				eventSummary = fmt.Sprintf("%d' %s\n%s\nIN\n%s\nOUT\n%s\n", eventTime, eventDetail, teamName, playerName, assistName)
			} else if eventType == "Goal" {
				eventSummary = fmt.Sprintf("%d' GOAL!!!\n%s\nPlayer: %s\nAssist: %s\n%s\n", eventTime, teamName, playerName, assistName, eventDetail)
			} else if eventType == "Var" {
				eventSummary = fmt.Sprintf("%d' %s\n%s\n%s\n%s\n", eventTime, eventType, teamName, playerName, eventDetail)
			}

			eventsArr = append(eventsArr, eventSummary)
		}

		fmt.Println(matchDisplay)
		for _, event := range eventsArr {
			fmt.Println(event)
		}
	},
}

func init() {
	rootCmd.AddCommand(liveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// liveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// liveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}