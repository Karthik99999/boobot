package mmr

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"boobot/structs"
)

func GetHlData(id string) (*structs.HlorenziBoard, string) {
	url := "https://gb.hlorenzi.com/api/v1/graphql"

	var jsonStr = []byte(`{
		team(teamId: "` + id + `") {
			kind,
			name,
			tiers {
				name,
				lowerBound,
				color
			}
			playerCount,
			players {
				name,
				ranking,
				maxRanking,
				minRanking,
				wins,
				losses,
				playedMatchCount,
				firstActivityDate,
				lastActivityDate,
				rating,
				ratingGain,
				maxRating,
				minRating,
				maxRatingGain,
				maxRatingLoss,
				points,
				maxPointsGain
			}
		}
	}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, "An error occured during the request. Please try again in a moment."
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, "An error occured during the request. Please try again in a moment."
	}
	var board *structs.HlorenziBoard
	_ = json.Unmarshal(body, &board)
	board.Data.Team.Url = "https://gb.hlorenzi.com/reg/" + id

	return board, ""
}
