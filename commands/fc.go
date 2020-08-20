package commands

import (
	"fmt"
	"log"
	"strings"

	"boobot/structs"

	"github.com/bwmarrin/discordgo"
	bolt "go.etcd.io/bbolt"
)

// Add command to list of commands
func init() {
	cmd := FC()
	Commands = append(Commands, cmd)
	fmt.Printf("loaded command: %s\n", cmd.Name)
}

// Initialize command
func FC() Command {
	cmd := Command{}
	cmd.Name = "fc"
	cmd.Run = runFC
	return cmd
}

// Function to run when command is used
func runFC(s *discordgo.Session, message *discordgo.MessageCreate, args []string, settings structs.GuildSettings) {
	// Open the database
	db, err := bolt.Open("db/fc.db", 0600, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(message.GuildID))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	}); err != nil {
		log.Println(err)
		return
	}

	if err := db.View(func(tx *bolt.Tx) error {
		var user *discordgo.User
		if len(args) > 0 {
			guild, _ := s.Guild(message.GuildID)
			for _, m := range guild.Members {
				if strings.ToLower(m.Nick) == strings.ToLower(strings.Join(args, " ")) {
					user = m.User
				}
			}
		} else {
			user = message.Author
		}
		if user == nil {
			s.ChannelMessageSend(message.ChannelID, "Couldn't find a user by that nickname.")
			return nil
		}
		b := tx.Bucket([]byte(message.GuildID))
		v := b.Get([]byte(user.ID))
		if string(v) == "" {
			if len(args) > 0 {
				s.ChannelMessageSend(message.ChannelID, "Looks like this user doesn't have an fc set.")
			} else {
				s.ChannelMessageSend(message.ChannelID, "Looks like you don't have an fc set. Use the `setfc` command to set one.")
			}
			return nil
		}
		s.ChannelMessageSend(message.ChannelID, string(v))
		return nil
	}); err != nil {
		log.Println(err)
		return
	}
}
