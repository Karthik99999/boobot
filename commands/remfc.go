package commands

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	bolt "go.etcd.io/bbolt"
)

func init() {
	CommandHandlers["remfc"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// Open the database
		db, err := bolt.Open("db/fc.db", 0600, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer db.Close()

		if err := db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(i.GuildID))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			return nil
		}); err != nil {
			log.Println(err)
			return
		}

		if err := db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(i.GuildID))
			err := b.Delete([]byte(i.Member.User.ID))
			if err != nil {
				log.Println(err)
				respond(s, i, "There was an error deleting your friend code. Please try again in a moment.", true)
			} else {
				respond(s, i, "Your friend code has been deleted.", true)
			}
			return err
		}); err != nil {
			log.Println(err)
			return
		}
	}
}
