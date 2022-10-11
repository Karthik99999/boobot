package commands

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"boobot/mmr"
	"boobot/structs"
	"boobot/utils"

	"github.com/bwmarrin/discordgo"
)

func init() {
	CommandHandlers["stats"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		args := utils.GetArgs(i.ApplicationCommandData().Options)
		settings := structs.GetSettings(i.GuildID)
		if settings.GameBoards1 != "" {
			// create embed
			guild, _ := s.Guild(i.GuildID)
			var embed *discordgo.MessageEmbed = new(discordgo.MessageEmbed)
			embed.Author = &discordgo.MessageEmbedAuthor{
				Name:    "Stats",
				IconURL: guild.IconURL(),
			}
			// tr will either be RT or CT. Irrelevant if there is no second leaderboard
			tr := args[0]
			name := args[1]
			var leaderboard *structs.HlorenziBoard
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

			var hlPlayer structs.HlPlayer
			for _, player := range leaderboard.Data.Team.Players {
				if strings.ToLower(player.Name) == strings.ToLower(name) {
					hlPlayer = player
				}
			}

			// Add each stat as a field
			addStatField := func(statName, statValue string) {
				field := &discordgo.MessageEmbedField{
					Name:   statName,
					Value:  statValue,
					Inline: true,
				}
				embed.Fields = append(embed.Fields, field)
			}
			// Find tier based on rating
			getTier := func(rating float64) structs.HlTiers {
				tiers := leaderboard.Data.Team.Tiers
				if len(tiers) < 1 {
					return structs.HlTiers{
						Name:  "???",
						Color: "#ffffff",
					}
				}
				tierData := tiers[len(tiers)-1]
				for i, tier := range tiers {
					if tier.LowerBound > int(rating) {
						if i == 0 {
							tierData = tiers[i]
						} else {
							tierData = tiers[i-1]
						}
						break
					}
				}
				return tierData
			}
			// populate embed
			tierColor, _ := strconv.ParseInt(strings.ReplaceAll(getTier(math.Floor(hlPlayer.Rating)).Color, "#", ""), 16, 64)
			embedColor := strconv.FormatInt(tierColor, 10)
			embed.Color, _ = strconv.Atoi(embedColor)
			if hlPlayer.Name == "" {
				embed.Footer = &discordgo.MessageEmbedFooter{
					Text: "The specified player wasn't found. Check your input for errors.",
				}
			} else {
				embed.Description = fmt.Sprintf("[%s](%s)", hlPlayer.Name, leaderboard.Data.Team.Url+"/player/"+strings.ReplaceAll(hlPlayer.Name, " ", "%20"))
				addStatField("Rank", "#"+strconv.Itoa(hlPlayer.Ranking+1))
				addStatField("Tier", getTier(math.Floor(hlPlayer.Rating)).Name)
				addStatField("Matches", strconv.Itoa(hlPlayer.PlayedMatchCount))
				addStatField("Rating", strconv.Itoa(int(math.Floor(hlPlayer.Rating))))
				addStatField("Max Rating", strconv.Itoa(int(math.Floor(hlPlayer.MaxRating))))
				addStatField("Min Rating", strconv.Itoa(int(math.Floor(hlPlayer.MinRating))))
				addStatField("Wins", strconv.Itoa(hlPlayer.Wins))
				addStatField("Losses", strconv.Itoa(hlPlayer.Losses))
				addStatField("Win Ratio", fmt.Sprintf("%.1f%%", (float64(hlPlayer.Wins)/float64(hlPlayer.PlayedMatchCount))*100))
				addStatField("Max Rating Gain", fmt.Sprintf("%+d", hlPlayer.MaxRatingGain))
				addStatField("Max Rating Loss", strconv.Itoa(hlPlayer.MaxRatingLoss))
				addStatField("Max Points", strconv.Itoa(hlPlayer.MaxPointsGain))
				addStatField("Avg Points", fmt.Sprintf("%.1f", (float64(hlPlayer.Points)/float64(hlPlayer.PlayedMatchCount))))
				addStatField("Best Rank", "#"+strconv.Itoa(hlPlayer.MinRanking+1))
				addStatField("Worst Rank", "#"+strconv.Itoa(hlPlayer.MaxRanking+1))
			}

			respondEmbed(s, i, embed, false)
		} else {
			if settings.Spreadsheet1 == "" || settings.SheetName == "" || settings.PlayerIndex == "" || settings.StatIndexes == "" {
				respond(s, i, "One or more settings required to use this command have not been set. Tell a moderator to use the `/changesetting` command to set them.", false)
				return
			}
			// Create embed
			guild, _ := s.Guild(i.GuildID)
			var embed *discordgo.MessageEmbed = new(discordgo.MessageEmbed)
			embed.Author = &discordgo.MessageEmbedAuthor{
				Name:    "Stats",
				IconURL: guild.IconURL(),
			}
			tr := args[0]
			name := args[1]
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
			// Send error message if there is one
			if leaderboard == nil {
				respond(s, i, errMsg, false)
				return
			}
			// Convert indexes from settings to ints
			reg := regexp.MustCompile(`\s`)
			playerIndex, err := strconv.Atoi(settings.PlayerIndex)
			if err != nil || playerIndex >= len(leaderboard[0]) {
				respond(s, i, "There was an error retrieving data from the leaderboard. Make sure the indexes for the player and stat columns are correct in the guild settings.", false)
				return
			}
			statIndexes := strings.Split(reg.ReplaceAllString(settings.StatIndexes, ""), ",")
			for _, statIndex := range statIndexes {
				index, err := strconv.Atoi(statIndex)
				if err != nil || index >= len(leaderboard[0]) {
					respond(s, i, "There was an error retrieving data from the leaderboard. Make sure the indexes for the player and stat columns are correct in the guild settings.", false)
					return
				}
			}
			// Add each stat as a field
			addStatField := func(statName, statValue string) {
				field := &discordgo.MessageEmbedField{
					Name:   statName,
					Value:  statValue,
					Inline: true,
				}
				embed.Fields = append(embed.Fields, field)
			}
			// Loop over leaderboard rows
			for _, row := range leaderboard {
				if len(row) <= playerIndex {
					continue
				}
				if strings.ToLower(row[playerIndex].(string)) == name {
					embed.Description = row[playerIndex].(string)
					for _, index := range statIndexes {
						i, _ := strconv.Atoi(index)
						addStatField(leaderboard[0][i].(string), row[i].(string))
					}
					break
				}
			}
			// Let the user know if the player wasn't found.
			if embed.Description == "" {
				embed.Footer = &discordgo.MessageEmbedFooter{
					Text: "The specified player wasn't found. Check your input for errors.",
				}
			}
			respondEmbed(s, i, embed, false)
		}
	}
}
