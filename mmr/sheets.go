package mmr

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/sheets/v4"
)

type credentials struct {
	Email        string `json:"client_email"`
	PrivateKeyID string `json:"private_key_id"`
	PrivateKey   string `json:"private_key"`
	TokenURL     string `json:"token_uri"`
}

func GetSSData(spreadsheetId, sheetName string) ([][]interface{}, string) {
	file, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Printf("Error reading credentials file: %v", err)
	}
	creds := credentials{}
	_ = json.Unmarshal([]byte(file), &creds)

	conf := &jwt.Config{
		Email:        creds.Email,
		PrivateKeyID: creds.PrivateKeyID,
		PrivateKey:   []byte(creds.PrivateKey),
		TokenURL:     creds.TokenURL,
		Scopes:       []string{"https://www.googleapis.com/auth/spreadsheets.readonly"},
	}

	client := conf.Client(oauth2.NoContext)

	srv, err := sheets.New(client)
	if err != nil {
		log.Printf("Unable to retrieve Sheets client: %v", err)
	}

	// Define the Sheet Name and fields to select
	readRange := sheetName + "!A:Z"

	// Pull the data from the sheet
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Printf("Unable to retrieve data from sheet: %v", err)
		return nil, "There was an error retrieving data from the leaderboard. Make sure the spreadsheet ID and sheet name are correct in the guild settings."
	}

	return resp.Values, ""
}
