package commands

import (
	"boobot/utils"

	"github.com/bwmarrin/discordgo"
)

func init() {
	CommandHandlers["choose"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		options := utils.GetArgs(i.ApplicationCommandData().Options)
		respond(s, i, utils.RandomVal(options).(string), true)
	}
}
