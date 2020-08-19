package commands

import (
	"fmt"

	"boobot/structs"
	"boobot/utils"

	"github.com/bwmarrin/discordgo"
)

// Add command to list of commands
func init() {
	cmd := Choose()
	Commands = append(Commands, cmd)
	fmt.Printf("loaded command: %s\n", cmd.Name)
}

// Initialize command
func Choose() Command {
	cmd := Command{}
	cmd.Name = "choose"
	cmd.Run = runChoose
	return cmd
}

// Function to run when command is used
func runChoose(s *discordgo.Session, message *discordgo.MessageCreate, args []string, settings structs.GuildSettings) {
	if settings.DisableChoose != "true" {
		s.ChannelMessageSend(message.ChannelID, utils.RandomVal(args).(string))
	}
}
