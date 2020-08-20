package commands

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"boobot/mmr"
	"boobot/structs"
	"boobot/utils"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/Iwark/spreadsheet.v2"
)

// Add command to list of commands
func init() {
	cmd := Stats()
	Commands = append(Commands, cmd)
	fmt.Printf("loaded command: %s\n", cmd.Name)
}

// Initialize command
func Stats() Command {
	cmd := Command{}
	cmd.Name = "stats"
	cmd.Run = runStats
	return cmd
}

// Function to run when command is used
func runStats(s *discordgo.Session, message *discordgo.MessageCreate, args []string, settings structs.GuildSettings) {
	var tr string
	if len(args) > 0 {
		tr = strings.ToLower(args[0])
	}
	// mkblounge
	if message.GuildID == "387347467332485122" || message.GuildID == "513093856338640916" {
		if tr == "rt" || tr == "ct" {
			var players []*structs.Player
			if len(args) < 2 {
				players = mmr.GetPlayers(tr, []string{message.Member.Nick})
			} else {
				players = mmr.GetPlayers(tr, []string{strings.Join(args, " ")[3:]})
			}
			// Create embed
			guild, _ := s.Guild(message.GuildID)
			var embed *discordgo.MessageEmbed = new(discordgo.MessageEmbed)
			embed.Author = &discordgo.MessageEmbedAuthor{
				Name:    "Stats",
				IconURL: guild.IconURL(),
			}
			// Stop if the player wasn't found
			if len(players) == 0 {
				embed.Footer = &discordgo.MessageEmbedFooter{
					Text: "The specified player wasn't found. Check your input for errors.",
				}
				s.ChannelMessageSendEmbed(message.ChannelID, embed)
				return
			}
			player := players[0]
			embed.Description = fmt.Sprintf("[%s](%s)", player.Name, player.URL)
			// Add each stat as a field
			addStatField := func(statName, statValue string) {
				field := &discordgo.MessageEmbedField{
					Name:   statName,
					Value:  statValue,
					Inline: true,
				}
				embed.Fields = append(embed.Fields, field)
			}
			// Add stats as fields
			addStatField("Rank", player.Ranking)
			addStatField("Percentile", utils.Ternary(player.TotalWars == 0, "-", utils.Nth(int(math.Floor(player.Percentile)))).(string))
			addStatField("MMR", strconv.Itoa(player.CurrentMmr))
			addStatField("Peak MMR", fmt.Sprintf("%v", utils.Ternary(player.TotalWars <= 10, "-", player.PeakMmr)))
			addStatField("Lowest MMR", fmt.Sprintf("%v", utils.Ternary(player.TotalWars <= 10, "-", player.LowestMmr)))
			addStatField("Events", strconv.Itoa(player.TotalWars))
			addStatField("Win %", fmt.Sprintf("%v", utils.Ternary(player.TotalWars == 0, "-", math.Round(player.WinPercentage*100))))
			addStatField("W-L", fmt.Sprintf("%v", utils.Ternary(player.TotalWars == 0, "-", fmt.Sprintf("%d-%d", player.Wins, player.Loss))))
			addStatField("W-L (Last 10)", fmt.Sprintf("%v", utils.Ternary(player.TotalWars == 0, "-", fmt.Sprintf("%d-%d", player.Wins10, player.Loss10))))
			addStatField("Gain/Loss (Last 10)", fmt.Sprintf("%v", utils.Ternary(player.TotalWars == 0, "-", fmt.Sprintf("%+d", player.Gainloss10Mmr))))
			addStatField("Max MMR Gain", fmt.Sprintf("%v", utils.Ternary(player.TotalWars == 0, "-", fmt.Sprintf("%+d", player.MaxGainMmr))))
			addStatField("Max MMR Loss", fmt.Sprintf("%v", utils.Ternary(player.TotalWars == 0, "-", player.MaxLossMmr)))
			addStatField("Avg. Score", fmt.Sprintf("%v", utils.Ternary(player.TotalWars == 0, "-", fmt.Sprintf("%.1f", player.AverageScore))))
			addStatField("Avg. Score (Last 10)", fmt.Sprintf("%v", utils.Ternary(player.TotalWars == 0, "-", fmt.Sprintf("%.1f", player.Average10Score))))
			addStatField("Std. Dev", fmt.Sprintf("%v", utils.Ternary(player.TotalWars == 0, "-", fmt.Sprintf("%.1f", player.StdScore))))
			addStatField("Std. Dev (Last 10)", fmt.Sprintf("%v", utils.Ternary(player.TotalWars == 0, "-", fmt.Sprintf("%.1f", player.Std10Score))))
			addStatField("Top Score", fmt.Sprintf("%v", utils.Ternary(player.TotalWars == 0, "-", player.TopScore)))
			// Add last updated date as footer
			date, _ := time.Parse("2006-01-02 15:04:05", player.UpdateDate)
			date = date.Add(-2 * time.Hour)
			embed.Footer = &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("Last updated on %s", date.In(time.Local).Format("January 02, 2006 at 3:04 PM (EDT)")),
			}
			s.ChannelMessageSendEmbed(message.ChannelID, embed)
		} else {
			s.ChannelMessageSend(message.ChannelID, "Please specify the leaderboard you would like to check.")
		}
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
		var leaderboard *spreadsheet.Sheet = nil
		var errMsg string
		if settings.Spreadsheet2 != "" {
			if len(args) > 1 {
				player = strings.Join(args[1:], " ")
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
		playerIndex, _ := strconv.Atoi(settings.PlayerIndex)
		statIndexes := strings.Split(reg.ReplaceAllString(settings.StatIndexes, ""), ",")
		for _, index := range statIndexes {
			i, err := strconv.Atoi(index)
			if err != nil {
				s.ChannelMessageSend(message.ChannelID, "There was an error retrieving data from the leaderboard. Make sure the indexes for the player and stat columns are correct in the guild settings using the `set` command.")
				return
			}
			if playerIndex >= len(leaderboard.Rows[0]) || i >= len(leaderboard.Rows[0]) {
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
		for _, row := range leaderboard.Rows {
			// Find player. Use the nickname if no name was specified
			if (settings.Spreadsheet2 == "" && len(args) > 0) || (settings.Spreadsheet2 != "" && len(args) > 1) {
				if strings.ToLower(row[playerIndex].Value) == strings.ToLower(player) {
					embed.Description = row[playerIndex].Value
					// add field for each stat column specificed in guild settings
					for _, index := range statIndexes {
						i, _ := strconv.Atoi(index)
						if i != playerIndex {
							addStatField(leaderboard.Rows[0][i].Value, row[i].Value)
						}
					}
				}
			} else {
				if strings.ToLower(row[playerIndex].Value) == strings.ToLower(message.Member.Nick) {
					embed.Description = row[playerIndex].Value
					for _, index := range statIndexes {
						i, _ := strconv.Atoi(index)
						addStatField(leaderboard.Rows[0][i].Value, row[i].Value)
					}
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
