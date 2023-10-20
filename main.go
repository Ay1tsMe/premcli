/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import "premcli/cmd"

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

func main() {
	cmd.Execute()
}
