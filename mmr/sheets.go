package mmr

import (
	"log"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/sheets/v4"
)

func GetSSData(spreadsheetId, sheetName string) ([][]interface{}, string) {
	conf := &jwt.Config{
		Email:        "leaderboards@boobot.iam.gserviceaccount.com",
		PrivateKeyID: "3943a01edcff60c33022542a576b487a60c223af",
		PrivateKey:   []byte("-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDKFen24oTX+tll\nqoRIgK7XUT1fRzkemq3DLFJ7gMtxcDrYxEfsxTNE6HhT45HLC3oALkQkf3YIw/D7\nhveerXginbVCSzqBxNQ8doL3pt4zQDnPvKzCXvf0lnr1OObEhEiYcEDr9DQevEA6\nH0JQr0quhvgvVaiMXplFAWXAn6FC+AQ2QJjDbXNBOsfUpJu/8WtfCariyEw9Dh34\nbnGeYIFheQAttfV/BUONbSQMhrB6Y01zCgEK0PZ69/iT0Zd2MF1plcZTiGoeE62E\ncfDE/paJ2IWO1/i05MsixeG8z/6Xx2ZNnhxA28FgUnSXc3dfiUAtcbhUqw9GRXRR\nSA8btkiTAgMBAAECggEAM4CzQJsJbTv+tOTo0suNA84uHH675XtZZqEAon44G0CV\nltIrXIIDp3+xzvt0GDHkFXC1KDId5Gz/mTMUH6opMHVOEUe38QO3bXNsvG4YOiqX\nsURuKRloCztgueeXFKV8FPGi8h+qOt0SZ125GnQaTfGTBglILAIeANKy2o00XasB\nAigoYqhKpgvEemHyyhnI4F9e0+S5nzDL0JAO3tnKkB6QU6XcWe0T6KY1IGIbYP7M\n5Qy8SQ7sUnYUvMgp9FO9Y94MOqKR7krjy8lzHtMEFJENkcUDg3eB/7QCOxI5wKRJ\nziB4f5MlX3uf7SbngAcGILp9yJnFV7dikL6i5xCPEQKBgQDmFVApVc8t3bUp2QQb\nlbtuhoHtW9sWjhe2i3JJozu6RiGygEScAvIRmuQs41jnhnKgYo5FoCdaIWkR6xLj\nU8hcnU53TOsgfuovKaQBXu/9uUkBDxPxjAKHlssrb/8zRBqCujhDqHtWhRZSCgpF\nMeGXc3cBNTpqnlusNb66dvzfUQKBgQDg2UJ86zh+qSdspJIJdMadMenuFD+zGyBy\nLgf1gq2iYHyv2ga6GhsHnkfMXckNNdf8YGGYJvuW21tWeNl0EUZ6wuRFgUO3HskV\nyIEGg+twKoVK12mEdWdfQzTXSFspb9x1D7gM/2kV7V/fEiFi70nui4N6dvD4Qep3\nkdHhA1yYowKBgHOegDLVWRAeWlxWHpdSDecDlqTVROo3qzjjKCJS8b+wYFyX0mJn\npIcuQ70+3b0ytcVc4UuhqETFh0wmyc4MmyHXNsgCkiE5Ras/jJfXwlfI1SPAFPCL\nv/Ws1BnW5PI5Je1NcNqm/pvCsy20t+Z/o3J85m9n9RwAyeZm95oyEu6RAoGBAJdf\njNR+oz0acjFBJhP5qxD/Hocq2Kui0pgsBy2w+WZ84NSeyrKViqb5V0rtxMIBAtSk\nqm99pxkrunUfzP3H95QECxwD52ur7SKeJscVHvcXmT8GgMItLBfFLhjVXJIr/dZN\na9JMTRn3AfkywolRoYtYH37d/mKUNd6jrBF5auhTAoGAcdYhQb/Z/93bf+PzMbko\nUcypGOzjWOVnvIb17GNr3ldz41kZC2hIJY4pr1A5FjBK8AaiReQSSjy/GVI7jf7g\nN2oqfv2AY+0YOD86N54ZwDguVrr34BbenEXIiVL1RmUAS6Wrv4yuogP7jQ01MPcj\nP4YOCZjSM8ZTM5c9vDlx0NM=\n-----END PRIVATE KEY-----\n"),
		TokenURL:     "https://oauth2.googleapis.com/token",
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
		return nil, "There was an error retrieving data from the leaderboard. Make sure the spreadsheet ID and sheet name are correct in the guild settings using the `set` command."
	}

	return resp.Values, ""
}
