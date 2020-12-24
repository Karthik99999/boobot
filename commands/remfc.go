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
	cmd := RemFC()
	Commands = append(Commands, cmd)
	fmt.Printf("loaded command: %s\n", cmd.Name)
}

// Initialize command
func RemFC() Command {
	cmd := Command{}
	cmd.Name = "remfc"
	cmd.Run = runRemFC
	return cmd
}

// Function to run when command is used
func runRemFC(s *discordgo.Session, message *discordgo.MessageCreate, args []string, settings structs.GuildSettings) {
	if strings.ToLower(settings.DisableFC) == "true" {
		return
	}
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

	if err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(message.GuildID))
		err := b.Delete([]byte(message.Author.ID))
		s.ChannelMessageSend(message.ChannelID, "Your FC has been deleted.")
		return err
	}); err != nil {
		log.Println(err)
		return
	}
}
