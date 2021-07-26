package commands

import (
	"boobot/structs"
	"fmt"
	"runtime/debug"

	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Name        string
	Aliases     []string
	Run         func(*discordgo.Session, *discordgo.MessageCreate, []string, structs.GuildSettings)
	Description string
	Cooldown    int
}

// Map of commands that will be added to on startup
var Commands map[string]Command

// Map of command aliases
var Aliases map[string]string

// Loads command into Commands map
func initCommand(cmd Command) {
	if Commands == nil {
		Commands = make(map[string]Command)
	}
	if Aliases == nil {
		Aliases = make(map[string]string)
	}
	Commands[cmd.Name] = cmd
	for _, a := range cmd.Aliases {
		Aliases[a] = cmd.Name
	}
	fmt.Printf("loaded command: %s\n", cmd.Name)
}

// Prevent bot from crashing from a panic
func recoverPanic(s *discordgo.Session, message *discordgo.MessageCreate) {
	if r := recover(); r != nil {
		fmt.Println("RECOVERED PANIC:", r)
		dmChannel, err := s.UserChannelCreate("397514708736802816")
		if err != nil {
			fmt.Println("Error sending panic stack trace DM:\n", err)
		} else {
			msg := fmt.Sprintf("RECOVERED PANIC:\n```%v\n%s```", r, string(debug.Stack()))
			s.ChannelMessageSend(dmChannel.ID, msg)
		}
	}
}
