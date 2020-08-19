package commands

import (
	"boobot/structs"
	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Name        string
	Aliases     []string
	Run         func(*discordgo.Session, *discordgo.MessageCreate, []string, structs.GuildSettings)
	Description string
	Cooldown    int
}

// List of commands that will be added to on startup
var Commands []Command
