package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"boobot/commands"
	"boobot/structs"
	"boobot/utils"

	"github.com/bwmarrin/discordgo"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	go http.ListenAndServe(":8080", nil)
	// Create a new Discord bot session using the token from a command line flag
	token := flag.String("token", "", "token to use to login to discord.")
	flag.Parse()
	if *token == "" {
		fmt.Println("A token wasn't provided! Use the -token flag to set one.")
		return
	}
	bot, err := discordgo.New("Bot " + *token)
	if err != nil {
		log.Fatal(err)
	}

	// register events
	bot.AddHandler(ready)
	bot.AddHandler(guildCreate)
	bot.AddHandler(guildDelete)
	bot.AddHandler(messageCreate)

	bot.Identify.Intents = nil

	err = bot.Open()

	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	// Cleanly close down the Discord session.
	bot.Close()
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	guildCount := fmt.Sprintf("%d servers", len(s.State.Guilds))
	usd := discordgo.UpdateStatusData{
		Game: &discordgo.Game{
			Name: "boobot.glitch.me | " + guildCount,
			Type: discordgo.GameTypeWatching,
		},
		Status: "dnd",
	}
	s.UpdateStatusComplex(usd)
	fmt.Println("logged in as user " + s.State.User.String())
}

func guildCreate(s *discordgo.Session, guild *discordgo.GuildCreate) {
	guildCount := fmt.Sprintf("%d servers", len(s.State.Guilds))
	usd := discordgo.UpdateStatusData{
		Game: &discordgo.Game{
			Name: "boobot.glitch.me | " + guildCount,
			Type: discordgo.GameTypeWatching,
		},
		Status: "dnd",
	}
	s.UpdateStatusComplex(usd)
}

func guildDelete(s *discordgo.Session, guild *discordgo.GuildDelete) {
	guildCount := fmt.Sprintf("%d servers", len(s.State.Guilds))
	usd := discordgo.UpdateStatusData{
		Game: &discordgo.Game{
			Name: "boobot.glitch.me | " + guildCount,
			Type: discordgo.GameTypeWatching,
		},
		Status: "dnd",
	}
	s.UpdateStatusComplex(usd)
}

func messageCreate(s *discordgo.Session, message *discordgo.MessageCreate) {
	// Don't reply to bots
	if message.Author.Bot {
		return
	}

	guildSettings := structs.GetSettings(message.GuildID)
	prefix := guildSettings.Prefix

	// Don't execute commands if they don't start with the prefix
	if strings.TrimPrefix(message.Content, prefix) == message.Content {
		return
	}

	args := strings.Fields(strings.TrimPrefix(message.Content, prefix))
	if len(args) <= 0 {
		return
	}
	command := strings.ToLower(args[0])
	args = args[1:]

	for _, c := range commands.Commands {
		if command == c.Name || utils.Contains(c.Aliases, command) {
			c.Run(s, message, args, guildSettings)
			guild, _ := s.Guild(message.GuildID)
			fmt.Printf("%s (%s) used %s command in %s (%s)\n", message.Author.Username, message.Author.ID, command, guild.Name, message.GuildID)
			break
		}
	}
}
