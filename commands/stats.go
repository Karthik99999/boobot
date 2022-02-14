package commands

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"boobot/mmr"
	"boobot/structs"

	"github.com/bwmarrin/discordgo"
)

func init() {
	cmd := Command{}
	cmd.Name = "stats"
	cmd.Run = runStats
	initCommand(cmd)
}

// Function to run when command is used
func runStats(s *discordgo.Session, message *discordgo.MessageCreate, args []string, settings structs.GuildSettings) {
	defer recoverPanic(s)
	if strings.ToLower(settings.DisableMMR) == "true" {
		return
	}
	if settings.GameBoards1 != "" {
		// create embed
		guild, _ := s.Guild(message.GuildID)
		var embed *discordgo.MessageEmbed = new(discordgo.MessageEmbed)
		embed.Author = &discordgo.MessageEmbedAuthor{
			Name:    "Stats",
			IconURL: guild.IconURL(),
		}
		var tr string
		if len(args) > 0 {
			tr = strings.ToLower(args[0])
		}
		var player string
		var leaderboard *structs.HlorenziBoard
		var errMsg string
		if settings.GameBoards2 != "" {
			if len(args) > 1 {
				player = strings.Join(args[1:], " ")
			} else if message.Member.Nick != "" {
				player = message.Member.Nick
			} else {
				player = message.Author.Username
			}
			if tr == "rt" {
				leaderboard, errMsg = mmr.GetHlData(settings.GameBoards1)
			} else if tr == "ct" {
				leaderboard, errMsg = mmr.GetHlData(settings.GameBoards2)
			} else {
				s.ChannelMessageSend(message.ChannelID, "Please specify the leaderboard you would like to check.")
				return
			}
		} else {
			if len(args) > 0 {
				player = strings.Join(args, " ")
			} else if message.Member.Nick != "" {
				player = message.Member.Nick
			} else {
				player = message.Author.Username
			}
			leaderboard, errMsg = mmr.GetHlData(settings.GameBoards1)
		}
		// Send error message if there is one
		if leaderboard == nil {
			s.ChannelMessageSend(message.ChannelID, errMsg)
			return
		} else if leaderboard.Data.Team.Kind != "lounge" {
			s.ChannelMessageSend(message.ChannelID, "The game boards ID wasn't found, or it was not for a lounge. Change it in the guild settings using the `set` command.")
			return
		}

		var hlPlayer structs.HlPlayer
		// Loop over players
		for _, players := range leaderboard.Data.Team.Players {
			// Find player. Use the nickname if no name was specified
			if (settings.GameBoards2 == "" && len(args) > 0) || (settings.GameBoards2 != "" && len(args) > 1) {
				if strings.ToLower(players.Name) == strings.ToLower(player) {
					hlPlayer = players
				}
			} else {
				if strings.ToLower(players.Name) == strings.ToLower(message.Member.Nick) || strings.ToLower(players.Name) == strings.ToLower(message.Author.Username) {
					hlPlayer = players
				}
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

		s.ChannelMessageSendEmbed(message.ChannelID, embed)
	} else {
		if settings.Spreadsheet1 == "" || settings.SheetName == "" || settings.PlayerIndex == "" || settings.StatIndexes == "" {
			s.ChannelMessageSend(message.ChannelID, "One or more settings required to use this command have not been set. Tell a moderator to use the `set` command to set them.")
			return
		}
		// Create embed
		guild, _ := s.Guild(message.GuildID)
		var embed *discordgo.MessageEmbed = new(discordgo.MessageEmbed)
		embed.Author = &discordgo.MessageEmbedAuthor{
			Name:    "Stats",
			IconURL: guild.IconURL(),
		}
		var tr string
		if len(args) > 0 {
			tr = strings.ToLower(args[0])
		}
		var player string
		var leaderboard [][]interface{}
		var errMsg string
		if settings.Spreadsheet2 != "" {
			if len(args) > 1 {
				player = strings.Join(args[1:], " ")
			} else if message.Member.Nick != "" {
				player = message.Member.Nick
			} else {
				player = message.Author.Username
			}
			if tr == "rt" {
				leaderboard, errMsg = mmr.GetSSData(settings.Spreadsheet1, settings.SheetName)
			} else if tr == "ct" {
				leaderboard, errMsg = mmr.GetSSData(settings.Spreadsheet2, settings.SheetName)
			} else {
				s.ChannelMessageSend(message.ChannelID, "Please specify the leaderboard you would like to check.")
				return
			}
		} else {
			if len(args) > 0 {
				player = strings.Join(args, " ")
			} else if message.Member.Nick != "" {
				player = message.Member.Nick
			} else {
				player = message.Author.Username
			}
			leaderboard, errMsg = mmr.GetSSData(settings.Spreadsheet1, settings.SheetName)
		}
		// Send error message if there is one
		if leaderboard == nil {
			s.ChannelMessageSend(message.ChannelID, errMsg)
			return
		}
		// Convert indexes from settings to ints
		reg := regexp.MustCompile(`\s`)
		playerIndex, err := strconv.Atoi(settings.PlayerIndex)
		if err != nil || playerIndex >= len(leaderboard[0]) {
			s.ChannelMessageSend(message.ChannelID, "There was an error retrieving data from the leaderboard. Make sure the indexes for the player and stat columns are correct in the guild settings using the `set` command.")
			return
		}
		statIndexes := strings.Split(reg.ReplaceAllString(settings.StatIndexes, ""), ",")
		for _, index := range statIndexes {
			i, err := strconv.Atoi(index)
			if err != nil || i >= len(leaderboard[0]) {
				s.ChannelMessageSend(message.ChannelID, "There was an error retrieving data from the leaderboard. Make sure the indexes for the player and stat columns are correct in the guild settings using the `set` command.")
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
			// Find player. Use the nickname if no name was specified
			if (settings.Spreadsheet2 == "" && len(args) > 0) || (settings.Spreadsheet2 != "" && len(args) > 1) {
				if strings.ToLower(row[playerIndex].(string)) == strings.ToLower(player) {
					embed.Description = row[playerIndex].(string)
					// add field for each stat column specificed in guild settings
					for _, index := range statIndexes {
						i, _ := strconv.Atoi(index)
						if i != playerIndex {
							addStatField(leaderboard[0][i].(string), row[i].(string))
						}
					}
					break
				}
			} else {
				if strings.ToLower(row[playerIndex].(string)) == strings.ToLower(message.Member.Nick) || strings.ToLower(row[playerIndex].(string)) == strings.ToLower(message.Author.Username) {
					embed.Description = row[playerIndex].(string)
					for _, index := range statIndexes {
						i, _ := strconv.Atoi(index)
						addStatField(leaderboard[0][i].(string), row[i].(string))
					}
					break
				}
			}
		}
		// Let the user know if the player wasn't found.
		if embed.Description == "" {
			embed.Footer = &discordgo.MessageEmbedFooter{
				Text: "The specified player wasn't found. Check your input for errors.",
			}
		}
		s.ChannelMessageSendEmbed(message.ChannelID, embed)
	}
}
