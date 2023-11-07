/*
Copyright © 2023 Adam Wyatt
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	roundValue string
)

type ApiResponseFixture struct {
	Response []Match `json:"response"`
}

type Match struct {
	Fixture struct {
		ID     int
		Date   string
		Status struct {
			Short   string
			Elapsed int
		}
	}
	Teams struct {
		Home struct {
			Name string
		}
		Away struct {
			Name string
		}
	}
	Goals struct {
		Home int
		Away int
	}
}

type CurrentRound struct {
	Response []string `json:"response"`
}

func getCurrentRound(previous bool, next bool) error {
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
	} else {
		return fmt.Errorf("No round information found in the API response")
	}

	// Edit currentRound to if previous or next flag called
	// This is kind of stupid but it's faster than another API call
	if previous || next {
		parts := strings.Split(roundValue, " ")
		if len(parts) > 0 {
			lastPart := parts[len(parts)-1]
			roundNum, err := strconv.Atoi(lastPart)

			if err != nil {
				return fmt.Errorf("Error parsing round number: %v", err)
			}

			if roundNum > 1 {
				if previous {
					roundNum--

				} else {
					roundNum++
				}

				roundValue = strings.Join(parts[:len(parts)-1], " ") + fmt.Sprintf(" %d", roundNum)
			} else {
				return fmt.Errorf("There is no previous round availible")
			}
		} else {
			return fmt.Errorf("Unexpected round format: %v", roundValue)
		}
	}

	return nil
}

func buildFixturesURL() string {
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

func getFixtures() ([]Match, error) {
	url := buildFixturesURL()

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

	var responseData ApiResponseFixture
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

func getDateFromMatchDisplay(matchDisplay string) string {
	dateString := strings.Split(matchDisplay, "Date: ")[1][:21]
	return dateString
}

func sortFixtures(fixturesArr []string) []string {
	layout := "02 Jan 2006, 03:04 PM"

	sort.Slice(fixturesArr, func(i, j int) bool {
		iDateString := getDateFromMatchDisplay(fixturesArr[i])
		jDateString := getDateFromMatchDisplay(fixturesArr[j])

		iDate, _ := time.Parse(layout, iDateString)
		jDate, _ := time.Parse(layout, jDateString)

		return iDate.Before(jDate)
	})

	return fixturesArr
}

func extractTeams(match string) (string, string) {
	if strings.Contains(match, " vs. ") {
		// Handle the "[H] Sheffield Utd vs. Wolves [A]" format
		parts := strings.Split(match, " vs. ")
		homeTeam := strings.Trim(parts[0][4:], " ")
		awayTeam := strings.Trim(parts[1][:len(parts[1])-4], " ")
		return homeTeam, awayTeam
	} else {
		// Handle the "[H] Sheffield Utd 2 - 3 Wolves [A]" format
		parts := strings.Split(match, " - ")
		home := strings.Split(parts[0], " ")
		away := strings.Split(parts[1], " ")

		homeTeam := strings.Trim(home[len(home)-2][4:], " ") // Taking the second-last word after [H] prefix
		awayTeam := strings.Trim(away[1], " ")               // Taking the first word before [A] suffix

		return homeTeam, awayTeam
	}
}

// fixturesCmd represents the fixtures command
var fixturesCmd = &cobra.Command{
	Use:   "fixtures",
	Short: "Prints fixtures for current round",
	Long:  "Prints fixtures for current round and highlights in bold the fixture of your favourite team.",
	Run: func(cmd *cobra.Command, args []string) {
		// Get flags if called
		previousRound, _ := cmd.Flags().GetBool("previous")
		nextRound, _ := cmd.Flags().GetBool("next")

		err := GetConfig()
		if err != nil {
			fmt.Println("Error loading config:", err)
			return
		}

		err = getCurrentRound(previousRound, nextRound)
		if err != nil {
			fmt.Println("Error getting current round:", err)
			return
		}

		matches, err := getFixtures()
		if err != nil {
			fmt.Println("Error fetching and parsing:", err)
			return
		}

		color.Set(color.Underline)
		fmt.Println(roundValue)
		color.Unset()

		var fixturesArr []string

		for _, match := range matches {
			homeTeam := match.Teams.Home.Name
			homeScore := match.Goals.Home
			awayTeam := match.Teams.Away.Name
			awayScore := match.Goals.Away
			date := match.Fixture.Date
			fixtureID := match.Fixture.ID
			timeElapsed := match.Fixture.Status.Elapsed
			matchStatus := match.Fixture.Status.Short

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
			vsPadding := vsWidth - len("[H]") - len(homeTeam) - len(fmt.Sprint(vsWidth))

			matchDisplay := ""
			if matchStatus == "NS" {
				// Match hasn't started
				// matchDisplay = fmt.Sprintf("[H] %s vs. %s [A]\nDate: %s\nFixture ID: %d\n", homeTeam, awayTeam, userFriendlyTime, fixtureID)
				matchDisplay = fmt.Sprintf("Date: %s\n[H] %s%*s%s\n[A] %s%*s\nStatus: Game Hasn't Started.\nFixture ID: %d\n", userFriendlyTime, homeTeam, vsPadding, "", "vs.", awayTeam, awayPadding, "", fixtureID)
			} else {
				if matchStatus == "FT" {
					// Match has finished
					// matchDisplay = fmt.Sprintf("[H] %s %d - %d %s [A]\nTime Elapsed: %d\nGAME HAS FINISHED!\nDate: %s\nFixture ID: %d\n", homeTeam, homeScore, awayScore, awayTeam, timeElapsed, userFriendlyTime, fixtureID)
					matchDisplay = fmt.Sprintf("Date: %s\n[H] %s%*s%d\n[A] %s%*s%d\nStatus: Game Has Finished.\nFixture ID: %d\n", userFriendlyTime, homeTeam, homePadding, "", homeScore, awayTeam, awayPadding, "", awayScore, fixtureID)
				} else {
					// Match in progress
					// matchDisplay = fmt.Sprintf("[H] %s %d - %d %s [A]\nTime Elapsed: %d\nDate: %s\nFixture ID: %d\n", homeTeam, homeScore, awayScore, awayTeam, timeElapsed, userFriendlyTime, fixtureID)
					matchDisplay = fmt.Sprintf("Date: %s\n[H] %s%*s%d\n[A] %s%*s%d\nTime Elapsed: %d\nFixture ID: %d\n", userFriendlyTime, homeTeam, homePadding, "", homeScore, awayTeam, awayPadding, "", awayScore, timeElapsed, fixtureID)
				}
			}

			fixturesArr = append(fixturesArr, matchDisplay)
		}
		sortFixtures(fixturesArr)
		for _, fixture := range fixturesArr {
			if isFavTeam(fixture, favTeam) {
				color.Set(color.FgMagenta)
				fmt.Println(fixture)
				color.Unset()
				continue
			}
			fmt.Println(fixture)
		}

	},
}

func init() {
	rootCmd.AddCommand(fixturesCmd)

	fixturesCmd.PersistentFlags().BoolP("previous", "p", false, "Get fixtures for the previous round")
	fixturesCmd.PersistentFlags().BoolP("next", "n", false, "Get fixtures for the next round")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fixturesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this /command
	// is called directly, e.g.:
	// fixturesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
