package commands

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"boobot/structs"
	"boobot/utils"

	"github.com/bwmarrin/discordgo"
	bolt "go.etcd.io/bbolt"
)

// Add command to list of commands
func init() {
	cmd := SetFC()
	Commands = append(Commands, cmd)
	fmt.Printf("loaded command: %s\n", cmd.Name)
}

// Initialize command
func SetFC() Command {
	cmd := Command{}
	cmd.Name = "setfc"
	cmd.Run = runSetFC
	cmd.Aliases = []string{"addfc"}
	return cmd
}

// Function to run when command is used
func runSetFC(s *discordgo.Session, message *discordgo.MessageCreate, args []string, settings structs.GuildSettings) {
	defer recoverPanic()
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
		reg, _ := regexp.Compile("[^0-9]+")
		fc := reg.ReplaceAllString(strings.Join(args, ""), "")
		fc = utils.InsertNth(fc, 4)
		if len(fc) != 14 {
			s.ChannelMessageSend(message.ChannelID, "FCs must be 12 digits long.")
			return nil
		}
		b := tx.Bucket([]byte(message.GuildID))
		err := b.Put([]byte(message.Author.ID), []byte(fc))
		s.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Your FC has been set to: %s", fc))
		return err
	}); err != nil {
		log.Println(err)
		return
	}
}
