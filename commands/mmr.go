package commands

import (
	"fmt"
	"strconv"
	"strings"

	"boobot/mmr"
	"boobot/structs"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/Iwark/spreadsheet.v2"
)

// Add command to list of commands
func init() {
	cmd := MMR()
	Commands = append(Commands, cmd)
	fmt.Printf("loaded command: %s\n", cmd.Name)
}

// Initialize command
func MMR() Command {
	cmd := Command{}
	cmd.Name = "mmr"
	cmd.Aliases = []string{"elo"}
	cmd.Run = runMMR
	return cmd
}

// Function to run when command is used
func runMMR(s *discordgo.Session, message *discordgo.MessageCreate, args []string, settings structs.GuildSettings) {
	// get player names seperated by commas
	cArgs := strings.Split(strings.Join(args, " "), ",")
	for i, p := range cArgs {
		cArgs[i] = strings.TrimSpace(p)
	}
	if len(cArgs[0]) > 3 {
		cArgs[0] = cArgs[0][3:]
	}
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
				players = mmr.GetPlayers(tr, cArgs)
			}
			guild, _ := s.Guild(message.GuildID)
			var embed *discordgo.MessageEmbed = new(discordgo.MessageEmbed)
			embed.Author = &discordgo.MessageEmbedAuthor{
				Name:    "MMR",
				IconURL: guild.IconURL(),
			}
			// Add each player as a field
			for _, p := range players {
				field := &discordgo.MessageEmbedField{
					Name:   p.Name,
					Value:  fmt.Sprintf("[%s](%s)", strconv.Itoa(p.CurrentMmr), p.URL),
					Inline: true,
				}
				embed.Fields = append(embed.Fields, field)
			}
			// Let the user know how many players weren't found
			missingPlayers := len(cArgs) - len(players)
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
			s.ChannelMessageSendEmbed(message.ChannelID, embed)
		} else {
			s.ChannelMessageSend(message.ChannelID, "Please specify the leaderboard you would like to check.")
		}
	} else {
		if settings.Spreadsheet1 == "" || settings.SheetName == "" || settings.PlayerIndex == "" || settings.RatingIndex == "" {
			s.ChannelMessageSend(message.ChannelID, "One or more settings required to use this command have not been set. Tell a moderator to use the `set` command to set them.")
			return
		}
		// Create embed
		guild, _ := s.Guild(message.GuildID)
		var embed *discordgo.MessageEmbed = new(discordgo.MessageEmbed)
		embed.Author = &discordgo.MessageEmbedAuthor{
			Name:    settings.RatingName,
			IconURL: guild.IconURL(),
		}
		// get player names seperated by commas
		cArgs := strings.Split(strings.Join(args, " "), ",")
		for i, p := range cArgs {
			cArgs[i] = strings.TrimSpace(p)
		}
		if settings.Spreadsheet2 != "" && len(cArgs[0]) > 3 {
			cArgs[0] = cArgs[0][3:]
		}
		var tr string
		if len(args) > 0 {
			tr = strings.ToLower(args[0])
		}
		var leaderboard *spreadsheet.Sheet = nil
		var errMsg string
		if settings.Spreadsheet2 != "" {
			if tr == "rt" {
				leaderboard, errMsg = mmr.GetSSData(settings.Spreadsheet1, settings.SheetName)
			} else if tr == "ct" {
				leaderboard, errMsg = mmr.GetSSData(settings.Spreadsheet2, settings.SheetName)
			} else {
				s.ChannelMessageSend(message.ChannelID, "Please specify the leaderboard you would like to check.")
				return
			}
		} else {
			leaderboard, errMsg = mmr.GetSSData(settings.Spreadsheet1, settings.SheetName)
		}
		// Send error message if there is one
		if leaderboard == nil {
			s.ChannelMessageSend(message.ChannelID, errMsg)
			return
		}

		playerIndex, err := strconv.Atoi(settings.PlayerIndex)
		if err != nil {
			s.ChannelMessageSend(message.ChannelID, "There was an error retrieving data from the leaderboard. Make sure the indexes for the player and rating columns are correct in the guild settings using the `set` command.")
			return
		}
		ratingIndex, err := strconv.Atoi(settings.RatingIndex)
		if err != nil {
			s.ChannelMessageSend(message.ChannelID, "There was an error retrieving data from the leaderboard. Make sure the indexes for the player and rating columns are correct in the guild settings using the `set` command.")
			return
		}
		if playerIndex >= len(leaderboard.Rows[0]) || ratingIndex >= len(leaderboard.Rows[0]) {
			s.ChannelMessageSend(message.ChannelID, "There was an error retrieving data from the leaderboard. Make sure the indexes for the player and rating columns are correct in the guild settings using the `set` command.")
			return
		}
		// Loop over leaderboard rows
		for _, row := range leaderboard.Rows {
			// Find player. Use the nickname if no name was specified
			if (settings.Spreadsheet2 == "" && len(args) > 0) || (settings.Spreadsheet2 != "" && len(args) > 1) {
				for _, player := range cArgs {
					if strings.ToLower(row[playerIndex].Value) == strings.ToLower(player) {
						field := &discordgo.MessageEmbedField{
							Name:   row[playerIndex].Value,
							Value:  row[ratingIndex].Value,
							Inline: true,
						}
						embed.Fields = append(embed.Fields, field)
					}
				}
			} else {
				if strings.ToLower(row[playerIndex].Value) == strings.ToLower(message.Member.Nick) {
					field := &discordgo.MessageEmbedField{
						Name:   row[playerIndex].Value,
						Value:  row[ratingIndex].Value,
						Inline: true,
					}
					embed.Fields = append(embed.Fields, field)
				}
			}
		}
		// Let the user know how many players weren't found.
		missingPlayers := len(cArgs) - len(embed.Fields)
		if missingPlayers > 0 {
			if missingPlayers == 1 {
				embed.Footer = &discordgo.MessageEmbedFooter{
					Text: "A player wasn't found. Check your input for errors.",
				}
			} else {
				embed.Footer = &discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("%d players wasn't found. Check your input for errors.", missingPlayers),
				}
			}
		}
		s.ChannelMessageSendEmbed(message.ChannelID, embed)
	}
}
