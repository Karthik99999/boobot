package commands

import (
	"strings"

	"boobot/structs"
	"boobot/utils"

	"github.com/bwmarrin/discordgo"
)

func init() {
	cmd := Command{}
	cmd.Name = "8ball"
	cmd.Run = run8ball
	initCommand(cmd)
}

// Function to run when command is used
func run8ball(s *discordgo.Session, message *discordgo.MessageCreate, args []string, settings structs.GuildSettings) {
	defer recoverPanic(s, message)
	responses := []string{
		"It is certain",
		"It is decidedly so",
		"Without a doubt",
		"Yes – definitely",
		"You may rely on it",
		"As I see it, yes",
		"Most likely",
		"Outlook good",
		"Yes",
		"Signs point to yes",
		"Reply hazy, try again",
		"Ask again later",
		"Better not tell you now",
		"Cannot predict now",
		"Concentrate and ask again",
		"Don't count on it",
		"My reply is no",
		"My sources say no",
		"Outlook not so good",
		"Very doubtful",
	}
	if len(args) < 1 || !strings.Contains(strings.Join(args, " "), "?") {
		s.ChannelMessageSend(message.ChannelID, "You didn't ask a question!")
		return
	}
	s.ChannelMessageSend(message.ChannelID, utils.RandomVal(responses).(string))
}
