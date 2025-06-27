package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"boobot/commands"
	"boobot/utils"

	"github.com/bwmarrin/discordgo"
)

func ready(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateGameStatus(0, fmt.Sprintf("%d servers", len(s.State.Guilds)))
	fmt.Println("logged in as user " + s.State.User.String())
}

func guildCreate(s *discordgo.Session, guild *discordgo.GuildCreate) {
	s.UpdateGameStatus(0, fmt.Sprintf("%d servers", len(s.State.Guilds)))
}

func guildDelete(s *discordgo.Session, guild *discordgo.GuildDelete) {
	s.UpdateGameStatus(0, fmt.Sprintf("%d servers", len(s.State.Guilds)))
}

func commandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	// If a command handler crashes, don't kill the bot
	defer utils.RecoverPanic(s, i)

	data := i.ApplicationCommandData()
	if h, ok := commands.CommandHandlers[data.Name]; ok {
		h(s, i)
	}
}

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
	s, err := discordgo.New("Bot " + *token)
	if err != nil {
		log.Fatal(err)
	}

	// register events
	s.AddHandler(ready)
	s.AddHandler(guildCreate)
	s.AddHandler(guildDelete)
	s.AddHandler(commandHandler)

	s.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = s.Open()
	if err != nil {
		log.Println("Error opening Discord session: ", err)
	}

	fmt.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands.Commands))
	for i, v := range commands.Commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", v)
		if err != nil {
			log.Printf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer s.Close()

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	fmt.Println("Removing commands...")
	for _, v := range registeredCommands {
		err := s.ApplicationCommandDelete(s.State.User.ID, "", v.ID)
		if err != nil {
			log.Printf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
}
