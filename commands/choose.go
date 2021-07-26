package commands

import (
	"boobot/structs"
	"boobot/utils"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func init() {
	cmd := Command{}
	cmd.Name = "choose"
	cmd.Run = runChoose
	initCommand(cmd)
}

// Function to run when command is used
func runChoose(s *discordgo.Session, message *discordgo.MessageCreate, args []string, settings structs.GuildSettings) {
	defer recoverPanic(s)
	if strings.ToLower(settings.DisableChoose) == "true" {
		return
	}
	if len(args) < 2 {
		s.ChannelMessageSend(message.ChannelID, "You need to provide at least 2 arguments.")
		return
	}
	s.ChannelMessageSend(message.ChannelID, utils.RandomVal(args).(string))
}
