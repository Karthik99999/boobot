package mmr

import (
	"context"
	"log"

	"golang.org/x/oauth2/google"
	"gopkg.in/Iwark/spreadsheet.v2"
)

func GetSSData(id, sheetName string) (*spreadsheet.Sheet, string) {
	conf, err := google.JWTConfigFromJSON(Creds, spreadsheet.Scope)
	if err != nil {
		log.Printf("Unable to read json data: %v", err)
	}

	client := conf.Client(context.TODO())
	service := spreadsheet.NewServiceWithClient(client)

	spreadsheet, err := service.FetchSpreadsheet(id)
	if err != nil {
		log.Printf("Unable to retrieve data from sheet: %v", err)
		return nil, "There was an error retrieving data from the leaderboard. Make sure the spreadsheet ID is correct in the guild settings using the `set` command."
	}
	sheet, err := spreadsheet.SheetByTitle(sheetName)
	if err != nil {
		log.Printf("Unable to retrieve data from sheet: %v", err)
		return nil, "There was an error retrieving data from the leaderboard. Make sure the sheet name is correct in the guild settings using the `set` command."
	}
	return sheet, ""
}
