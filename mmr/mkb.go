package mmr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"boobot/structs"
)

func GetPlayers(tracks string, names []string) []*structs.Player {
	url := fmt.Sprintf("https://mariokartboards.com/lounge/json/player.php?type=%s&name=%s", tracks, strings.ReplaceAll(strings.Join(names, ","), " ", ""))
	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return nil
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return nil
	}

	var players []*structs.Player
	_ = json.Unmarshal(body, &players)

	return players
}
