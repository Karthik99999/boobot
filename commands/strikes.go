package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"boobot/mmr"
	"boobot/structs"

	"github.com/bwmarrin/discordgo"
)

// Add command to list of commands
func init() {
	cmd := Strikes()
	Commands = append(Commands, cmd)
	fmt.Printf("loaded command: %s\n", cmd.Name)
}

// Initialize command
func Strikes() Command {
	cmd := Command{}
	cmd.Name = "strikes"
	cmd.Run = runStrikes
	return cmd
}

// Function to run when command is used
func runStrikes(s *discordgo.Session, message *discordgo.MessageCreate, args []string, settings structs.GuildSettings) {
	defer recoverPanic()
	if strings.ToLower(settings.DisableMMR) == "true" {
		return
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
				if message.Member.Nick == "" {
					players = mmr.GetPlayers(tr, []string{message.Author.Username})
				} else {
					players = mmr.GetPlayers(tr, []string{message.Member.Nick})
				}
			} else {
				players = mmr.GetPlayers(tr, []string{strings.Join(args, " ")[3:]})
			}
			// Create embed
			guild, _ := s.Guild(message.GuildID)
			var embed *discordgo.MessageEmbed = new(discordgo.MessageEmbed)
			embed.Author = &discordgo.MessageEmbedAuthor{
				Name:    "Strikes",
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
			// Add each strike stats as a field
			addStrikeField := func(statName, statValue string) {
				field := &discordgo.MessageEmbedField{
					Name:   statName,
					Value:  statValue,
					Inline: true,
				}
				embed.Fields = append(embed.Fields, field)
			}
			// Add strike stats as fields
			addStrikeField("Strikes", strconv.Itoa(player.Strikes))
			addStrikeField("Limit", strconv.Itoa(player.TotalStrikes))
			addStrikeField("Penalties", strconv.Itoa(player.Penalties))

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

	}
}
