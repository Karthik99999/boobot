package commands

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	bolt "go.etcd.io/bbolt"
)

func init() {
	CommandHandlers["fc"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

		if err := db.View(func(tx *bolt.Tx) error {
			var msg string
			user := i.ApplicationCommandData().Options[0].UserValue(s)
			b := tx.Bucket([]byte(i.GuildID))
			v := b.Get([]byte(user.ID))
			if string(v) == "" {
				msg = fmt.Sprintf("%s's FC: Not set", user.Username)
			} else {
				msg = fmt.Sprintf("%s's FC: %s", user.Username, string(v))
			}
			respond(s, i, msg, false)
			return nil
		}); err != nil {
			log.Println(err)
			return
		}
	}
}
