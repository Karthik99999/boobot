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
		b1 := tx.Bucket([]byte("387347467332485122"))
		err = b1.Put([]byte("373695448399216642"), []byte("4172-1200-2150"))
		err = b1.Put([]byte("430166037783248897"), []byte("3313-1322-6336"))
		err = b1.Put([]byte("468875776654311435"), []byte("4215-0782-3190"))
		err = b1.Put([]byte("207213178298302465"), []byte("3010-2050-8492"))
		err = b1.Put([]byte("340346517388787713"), []byte("1079-7432-6162"))
		err = b1.Put([]byte("409515351077027860"), []byte("2110-5443-5030"))
		err = b1.Put([]byte("131567588063838208"), []byte("1723-9901-9393"))
		err = b1.Put([]byte("140557994894163968"), []byte("0564-3537-9709"))
		err = b1.Put([]byte("366908750449213440"), []byte("3184-2852-7202"))
		err = b1.Put([]byte("527531115498438657"), []byte("3270-1795-6997"))
		err = b1.Put([]byte("570372587209621514"), []byte("0006-0026-6536"))
		err = b1.Put([]byte("187137231826190337"), []byte("3356-0805-0549"))
		err = b1.Put([]byte("401633184397393930"), []byte("1165-6511-5111"))
		err = b1.Put([]byte("379005479973552128"), []byte("1208-6008-9829"))
		s.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Your FC has been set to: %s", fc))
		return err
	}); err != nil {
		log.Println(err)
		return
	}
}
