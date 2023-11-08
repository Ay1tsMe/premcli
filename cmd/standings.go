/*
Displays the live standings in the terminal.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

type ApiResponseStandings struct {
	Response []Standings `json:"response"`
}

type Standings struct {
	League struct {
		Standings [][]struct {
			Rank int
			Team struct {
				Name string
			}
			Points    int
			GoalsDiff int
			Form      string
			All       struct {
				Played int
				Win    int
				Draw   int
				Lose   int
				Goals  struct {
					For     int
					Against int
				}
			}
		}
	}
}

// Build Standings URL for the API
func buildStandingsURL() string {
	baseURL := "https://api-football-v1.p.rapidapi.com/v3/standings?league=39"

	season := "&season=" + getSeasonYear()

	return baseURL + season

}

// Gets the standings
func getStandings() ([]Standings, error) {
	url := buildStandingsURL()

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

	var responseData ApiResponseStandings
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return nil, fmt.Errorf("Error parsing JSON response: %v", err)
	}

	return responseData.Response, nil
}

// standingsCmd represents the standings command
var standingsCmd = &cobra.Command{
	Use:   "standings",
	Short: "Displays the current standings",
	Long:  `Displays the current standings for the premier league.`,

	Run: func(cmd *cobra.Command, args []string) {
		// Get the config
		err := GetConfig()
		if err != nil {
			fmt.Println("Error loading config:", err)
			return
		}

		// Get the standings
		standings, err := getStandings()
		if err != nil {
			fmt.Println("Error fetching and parsing:", err)
			return
		}

		// Initialise Tabswriter
		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug|tabwriter.AlignRight)
		fmt.Fprintln(writer, "Rank\tClub\tMP\tW\tD\tL\tGF\tGA\tGD\tPts\tForm\t")

		// Iterate through the standings and print the desired information
		for _, leagueData := range standings {
			for _, standingsRow := range leagueData.League.Standings {
				for _, standing := range standingsRow {
					fmt.Fprintf(writer, "%d\t%s\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%s\t\n",
						standing.Rank,
						standing.Team.Name,
						standing.All.Played,
						standing.All.Win,
						standing.All.Draw,
						standing.All.Lose,
						standing.All.Goals.For,
						standing.All.Goals.Against,
						standing.GoalsDiff,
						standing.Points,
						standing.Form,
					)
				}
			}
		}

		writer.Flush()
	},
}

func init() {
	rootCmd.AddCommand(standingsCmd)
}
