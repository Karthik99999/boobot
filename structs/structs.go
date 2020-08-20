package structs

import (
	"bytes"
	"encoding/gob"
	"log"
	"reflect"
	"strings"

	"github.com/syndtr/goleveldb/leveldb"
)

type GuildSettings struct {
	guildID       string
	Prefix        string
	DisableChoose string
	Spreadsheet1  string
	Spreadsheet2  string
	SheetName     string
	RatingName    string
	PlayerIndex   string
	RatingIndex   string
	StatIndexes   string
}

func GetSettings(guildID string) GuildSettings {
	var gs GuildSettings
	// Open the database
	db, err := leveldb.OpenFile("db/settings", nil)
	if err != nil {
		log.Println(err)
		return GuildSettings{}
	}
	defer db.Close()

	v, _ := db.Get([]byte(guildID), nil)
	gs.Decode(v)
	gs.Default(guildID)
	return gs
}

// Set default settings. Should be called after initializing settings variable
func (gs *GuildSettings) Default(guildID string) {
	gs.guildID = guildID
	if gs.Prefix == "" {
		gs.Prefix = "$"
	}
	if gs.DisableChoose == "" {
		gs.DisableChoose = "false"
	}
	if gs.SheetName == "" {
		gs.SheetName = "Leaderboard"
	}
	if gs.RatingName == "" {
		gs.RatingName = "MMR"
	}
}

// Returns a byte array of the struct
func (gs GuildSettings) Encode() []byte {
	var encoded bytes.Buffer
	enc := gob.NewEncoder(&encoded)
	_ = enc.Encode(gs)
	return encoded.Bytes()
}

// Decodes a byte array into the struct
func (gs *GuildSettings) Decode(byteArr []byte) {
	decoded := bytes.NewBuffer(byteArr)
	dec := gob.NewDecoder(decoded)
	_ = dec.Decode(gs)
}

// View setting in guild. Used in set command to view settings based on command args
func (gs GuildSettings) View(key string) (string, bool) {
	var f = func(field string) bool {
		return strings.ToLower(field) == strings.ToLower(key)
	}
	if !reflect.Indirect(reflect.ValueOf(gs)).FieldByNameFunc(f).IsValid() {
		return "", false
	}
	return reflect.Indirect(reflect.ValueOf(gs)).FieldByNameFunc(f).String(), true
}

// Set setting in guild. Used in set command to change settings based on command args
func (gs *GuildSettings) Set(key, value string) bool {
	// Open the database
	db, err := leveldb.OpenFile("db/settings", nil)
	if err != nil {
		log.Println(err)
		return false
	}
	defer db.Close()

	var f = func(field string) bool {
		return strings.ToLower(field) == strings.ToLower(key)
	}
	if !reflect.ValueOf(gs).Elem().FieldByNameFunc(f).IsValid() {
		return false
	}
	reflect.ValueOf(gs).Elem().FieldByNameFunc(f).SetString(value)
	_ = db.Put([]byte(gs.guildID), gs.Encode(), nil)
	return true
}

// Reset guild setting to empty string (or default setting if it has one)
func (gs *GuildSettings) Reset(key string) bool {
	// Open the database
	db, err := leveldb.OpenFile("db/settings", nil)
	if err != nil {
		log.Println(err)
		return false
	}
	defer db.Close()

	var f = func(field string) bool {
		return strings.ToLower(field) == strings.ToLower(key)
	}
	if !reflect.ValueOf(gs).Elem().FieldByNameFunc(f).IsValid() {
		return false
	}
	reflect.ValueOf(gs).Elem().FieldByNameFunc(f).SetString("")
	_ = db.Put([]byte(gs.guildID), gs.Encode(), nil)
	return true
}

type Player struct {
	Pid             int     `json:"pid"`
	Name            string  `json:"name"`
	Strikes         int     `json:"strikes"`
	CurrentMmr      int     `json:"current_mmr"`
	PeakMmr         int     `json:"peak_mmr"`
	LowestMmr       int     `json:"lowest_mmr"`
	Wins            int     `json:"wins"`
	Loss            int     `json:"loss"`
	MaxGainMmr      int     `json:"max_gain_mmr"`
	MaxLossMmr      int     `json:"max_loss_mmr"`
	WinPercentage   float64 `json:"win_percentage"`
	Gainloss10Mmr   int     `json:"gainloss10_mmr"`
	Wins10          int     `json:"wins10"`
	Loss10          int     `json:"loss10"`
	Win10Percentage float64 `json:"win10_percentage"`
	WinStreak       int     `json:"win_streak"`
	TopScore        int     `json:"top_score"`
	AverageScore    float64 `json:"average_score"`
	Average10Score  float64 `json:"average10_score"`
	StdScore        float64 `json:"std_score"`
	Std10Score      float64 `json:"std10_score"`
	TotalWars       int     `json:"total_wars"`
	Penalties       int     `json:"penalties"`
	TotalStrikes    int     `json:"total_strikes"`
	Ranking         string  `json:"ranking"`
	Percentile      float64 `json:"percentile"`
	UpdateDate      string  `json:"update_date"`
	URL             string  `json:"url"`
}
