/*
Displays the live events timeline of a fixture given the fixtureID.
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

// Builds the API URL for retrieving fixtures
func buildFixtureByIDURL(fixtureID int) string {
	baseURL := "https://api-football-v1.p.rapidapi.com/v3/fixtures?"

	fixture := "id=" + strconv.Itoa(fixtureID)

	return baseURL + fixture
}

// Builds the API URL for retrieving events
func buildEventsURL(fixtureID int) string {
	baseURL := "https://api-football-v1.p.rapidapi.com/v3/fixtures/events?"

	fixture := "fixture=" + strconv.Itoa(fixtureID)

	return baseURL + fixture
}

// Gets the events given a fixture ID and parses the JSON
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

// Gets the fixture information given a fixture ID and parses the JSON
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

var liveCmd = &cobra.Command{
	Use:   "live <fixtureID>",
	Short: "Tracks the live events of a fixture",
	Long: `Tracks the live events of a fixture given a 'fixtureID'.

To obtain the 'fixtureID', use 'premcli fixtures' to display the fixtures with their appropriate 'fixtureID'.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get the config
		err := GetConfig()
		if err != nil {
			fmt.Println("Error loading config:", err)
			return
		}

		fixtureID, err := strconv.Atoi(args[0])
		if err != nil {
			return
		}

		// Get events
		events, err := getEvents(fixtureID)
		if err != nil {
			fmt.Println("Error fetching and parsing:", err)
			return
		}

		// Get fixture information
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

		// Loop through each event and store it in output array
		for _, event := range events {
			eventTime := event.Time.Elapsed
			eventExtraTime := event.Time.Extra
			teamName := event.Team.Name
			playerName := event.Player.Name
			assistName := event.Assist.Name
			eventType := event.Type
			eventDetail := event.Detail
			eventComment := event.Comments

			// Compute Time
			extraTimeStr := ""
			if eventExtraTime > 0 {
				extraTimeStr = fmt.Sprintf("+%d", eventExtraTime)
			}

			eventSummary := ""
			if eventType == "Card" {
				eventSummary = fmt.Sprintf("%d'%s %s\n%s\n%s\n%s\n", eventTime, extraTimeStr, eventDetail, teamName, playerName, eventComment)
			} else if eventType == "subst" {
				eventSummary = fmt.Sprintf("%d'%s %s\n%s\nIN\n%s\nOUT\n%s\n", eventTime, extraTimeStr, eventDetail, teamName, playerName, assistName)
			} else if eventType == "Goal" {
				eventSummary = fmt.Sprintf("%d'%s GOAL!!!\n%s\nPlayer: %s\nAssist: %s\n%s\n", eventTime, extraTimeStr, teamName, playerName, assistName, eventDetail)
			} else if eventType == "Var" {
				eventSummary = fmt.Sprintf("%d'%s %s\n%s\n%s\n%s\n", eventTime, extraTimeStr, eventType, teamName, playerName, eventDetail)
			}

			eventsArr = append(eventsArr, eventSummary)
		}

		// Display Fixture info and events
		fmt.Println(matchDisplay)
		for _, event := range eventsArr {
			fmt.Println(event)
		}
	},
}

func init() {
	rootCmd.AddCommand(liveCmd)

	liveCmd.Example = ` # Retrieve live events for fixture with ID 1234
premcli live 1234`
}
