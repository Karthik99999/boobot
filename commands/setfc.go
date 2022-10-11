package commands

import (
	"boobot/utils"
	"fmt"
	"log"
	"regexp"

	"github.com/bwmarrin/discordgo"
	bolt "go.etcd.io/bbolt"
)

func init() {
	CommandHandlers["setfc"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
			reg, _ := regexp.Compile("[^0-9]+")
			fc := reg.ReplaceAllString(i.ApplicationCommandData().Options[0].StringValue(), "")
			fc = utils.InsertNth(fc, 4)
			if len(fc) != 14 {
				respond(s, i, "FCs must be 12 digits long.", true)
				return nil
			}
			b := tx.Bucket([]byte(i.GuildID))
			err := b.Put([]byte(i.Member.User.ID), []byte(fc))
			if err != nil {
				log.Println(err)
				respond(s, i, "There was an error setting your friend code. Please try again in a moment.", true)
			} else {
				respond(s, i, fmt.Sprintf("Your FC has been set to: %s", fc), true)
			}
			return err
		}); err != nil {
			log.Println(err)
			return
		}
	}
}
