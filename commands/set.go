package commands

import (
	"fmt"
	"reflect"
	"strings"

	"boobot/structs"

	"github.com/bwmarrin/discordgo"
)

// Add command to list of commands
func init() {
	cmd := Set()
	Commands = append(Commands, cmd)
	fmt.Printf("loaded command: %s\n", cmd.Name)
}

// Initialize command
func Set() Command {
	cmd := Command{}
	cmd.Name = "set"
	cmd.Run = runSet
	cmd.Aliases = []string{"setting"}
	return cmd
}

// Function to run when command is used
func runSet(s *discordgo.Session, message *discordgo.MessageCreate, args []string, settings structs.GuildSettings) {
	p, _ := s.State.UserChannelPermissions(message.Author.ID, message.ChannelID)
	if p&discordgo.PermissionManageServer != discordgo.PermissionManageServer && message.Author.ID != "397514708736802816" {
		s.ChannelMessageSend(message.ChannelID, "You don't have the permission to use this command.")
		return
	}
	if len(args) > 0 {
		if strings.ToLower(args[0]) == "view" {
			if len(args) < 2 {
				s.ChannelMessageSend(message.ChannelID, "You didn't provide a guild setting to view!")
				return
			}
			setting, valid := settings.View(args[1])
			if !valid {
				s.ChannelMessageSend(message.ChannelID, "A guild setting with that name doesn't exist!")
				return
			}
			s.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Value of setting %s: %s", args[1], setting))
		} else if strings.ToLower(args[0]) == "edit" {
			if len(args) < 2 {
				s.ChannelMessageSend(message.ChannelID, "You didn't provide a guild setting to edit!")
				return
			}
			if len(args) < 3 {
				s.ChannelMessageSend(message.ChannelID, "You didn't provide a value to change the setting to!")
				return
			}
			valid := settings.Set(args[1], strings.Join(args[2:], " "))
			if !valid {
				s.ChannelMessageSend(message.ChannelID, "A guild setting with that name doesn't exist!")
				return
			}
			s.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Value of setting %s changed to: %s", args[1], strings.Join(args[2:], " ")))
		} else if strings.ToLower(args[0]) == "reset" {
			if len(args) < 2 {
				s.ChannelMessageSend(message.ChannelID, "You didn't provide a guild setting to reset!")
				return
			}
			valid := settings.Reset(args[1])
			if !valid {
				s.ChannelMessageSend(message.ChannelID, "A guild setting with that name doesn't exist!")
				return
			}
			s.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Guild setting %s reset to default value.", args[1]))
		}
	} else {
		msg := "```asciidoc"
		msg += "\n= Current guild settings ="

		k := reflect.Indirect(reflect.ValueOf(settings))
		typeOfT := k.Type()

		for i := 1; i < k.NumField(); i++ {
			f := k.Field(i)
			msg += fmt.Sprintf("\n%s :: %v", typeOfT.Field(i).Name, f.Interface())
		}
		msg += fmt.Sprintf("\n%sset edit <SettingName> <Value> to change a setting.", settings.Prefix)
		msg += fmt.Sprintf("\n%sset reset <SettingName> to reset a setting to it's default (blank if it doesn't have one).", settings.Prefix)
		msg += "\nGo to https://boobot.glitch.me/help.html for more info on how to setup these guild settings."
		msg += "```"
		s.ChannelMessageSend(message.ChannelID, msg)
	}
}
