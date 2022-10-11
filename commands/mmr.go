package commands

import (
	"boobot/mmr"
	"boobot/structs"
	"boobot/utils"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func init() {
	CommandHandlers["mmr"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		args := utils.GetArgs(i.ApplicationCommandData().Options)
		settings := structs.GetSettings(i.GuildID)
		if settings.GameBoards1 != "" {
			// create embed
			guild, _ := s.Guild(i.GuildID)
			var embed *discordgo.MessageEmbed = new(discordgo.MessageEmbed)
			embed.Author = &discordgo.MessageEmbedAuthor{
				Name:    settings.RatingName,
				IconURL: guild.IconURL(),
			}
			// tr will either be RT or CT. Irrelevant if there is no second leaderboard
			tr := args[0]
			names := args[1:]
			var leaderboard *structs.HlorenziBoard = nil
			var errMsg string
			if settings.GameBoards2 != "" {
				if tr == "rt" {
					leaderboard, errMsg = mmr.GetHlData(settings.GameBoards1)
				} else if tr == "ct" {
					leaderboard, errMsg = mmr.GetHlData(settings.GameBoards2)
				}
			} else {
				leaderboard, errMsg = mmr.GetHlData(settings.GameBoards1)
			}

			if leaderboard == nil {
				respond(s, i, errMsg, false)
				return
			}

			for _, player := range leaderboard.Data.Team.Players {
				for _, name := range names {
					if strings.ToLower(player.Name) == strings.ToLower(name) {
						field := &discordgo.MessageEmbedField{
							Name:   player.Name,
							Value:  fmt.Sprintf("[%s](%s)", strconv.Itoa(int(math.Floor(player.Rating))), leaderboard.Data.Team.Url+"/player/"+strings.ReplaceAll(player.Name, " ", "%20")),
							Inline: true,
						}
						embed.Fields = append(embed.Fields, field)
					}
				}
			}
			// Let the user know how many players weren't found.
			missingPlayers := len(names) - len(embed.Fields)
			if missingPlayers > 0 {
				if missingPlayers == 1 {
					embed.Footer = &discordgo.MessageEmbedFooter{
						Text: "A player wasn't found. Check your input for errors.",
					}
				} else {
					embed.Footer = &discordgo.MessageEmbedFooter{
						Text: fmt.Sprintf("%d players weren't found. Check your input for errors.", missingPlayers),
					}
				}
			}
			respondEmbed(s, i, embed, false)
		} else {
			if settings.Spreadsheet1 == "" || settings.SheetName == "" || settings.PlayerIndex == "" || settings.RatingIndex == "" {
				respond(s, i, "One or more settings required to use this command have not been set. Change them in the guild settings.", false)
				return
			}
			// Create embed
			guild, _ := s.Guild(i.GuildID)
			var embed *discordgo.MessageEmbed = new(discordgo.MessageEmbed)
			embed.Author = &discordgo.MessageEmbedAuthor{
				Name:    settings.RatingName,
				IconURL: guild.IconURL(),
			}
			// tr will either be rt or ct. Irrelevant if there is no second leaderboard
			tr := args[0]
			names := args[1:]
			var leaderboard [][]interface{}
			var errMsg string
			if settings.Spreadsheet2 != "" {
				if tr == "rt" {
					leaderboard, errMsg = mmr.GetSSData(settings.Spreadsheet1, settings.SheetName)
				} else if tr == "ct" {
					leaderboard, errMsg = mmr.GetSSData(settings.Spreadsheet2, settings.SheetName)
				}
			} else {
				leaderboard, errMsg = mmr.GetSSData(settings.Spreadsheet1, settings.SheetName)
			}

			if leaderboard == nil {
				respond(s, i, errMsg, true)
				return
			}

			playerIndex, err := strconv.Atoi(settings.PlayerIndex)
			if err != nil {
				respond(s, i, "There was an error retrieving data from the leaderboard. Make sure the indexes for the player and rating columns are correct in the guild settings.", false)
				return
			}
			ratingIndex, err := strconv.Atoi(settings.RatingIndex)
			if err != nil {
				respond(s, i, "There was an error retrieving data from the leaderboard. Make sure the indexes for the player and rating columns are correct in the guild settings.", false)
				return
			}
			if playerIndex >= len(leaderboard[0]) || ratingIndex >= len(leaderboard[0]) {
				respond(s, i, "There was an error retrieving data from the leaderboard. Make sure the indexes for the player and rating columns are correct in the guild settings.", false)
				return
			}

			for _, row := range leaderboard {
				if len(row) <= playerIndex || len(row) <= ratingIndex {
					continue
				}
				for _, name := range names {
					if strings.ToLower(row[playerIndex].(string)) == strings.ToLower(name) {
						field := &discordgo.MessageEmbedField{
							Name:   row[playerIndex].(string),
							Value:  row[ratingIndex].(string),
							Inline: true,
						}
						embed.Fields = append(embed.Fields, field)
					}
				}
			}
			// Let the user know how many players weren't found.
			missingPlayers := len(names) - len(embed.Fields)
			if missingPlayers > 0 {
				if missingPlayers == 1 {
					embed.Footer = &discordgo.MessageEmbedFooter{
						Text: "A player wasn't found. Check your input for errors.",
					}
				} else {
					embed.Footer = &discordgo.MessageEmbedFooter{
						Text: fmt.Sprintf("%d players weren't found. Check your input for errors.", missingPlayers),
					}
				}
			}
			respondEmbed(s, i, embed, false)
		}
	}
}
