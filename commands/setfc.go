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
		err = b1.Put([]byte("322520254712643586"), []byte("0690-5527-1295"))
		err = b1.Put([]byte("138838927275458561"), []byte("5288-8186-5284"))
		err = b1.Put([]byte("342098621262462980"), []byte("2067-5865-5496"))
		err = b1.Put([]byte("433353529655296011"), []byte("1895-7902-4602"))
		err = b1.Put([]byte("277498894760542208"), []byte("3227-2313-5370"))
		err = b1.Put([]byte("217694332730212352"), []byte("5114-5039-7324"))
		err = b1.Put([]byte("478575909230739459"), []byte("2582-9869-1012"))
		err = b1.Put([]byte("553976484826578980"), []byte("2110-5409-7480"))
		err = b1.Put([]byte("503787557188927509"), []byte("1938-7404-5086"))
		err = b1.Put([]byte("249868601249497089"), []byte("4000-3273-3612"))
		err = b1.Put([]byte("465269517988134922"), []byte("0864-9983-5463"))
		err = b1.Put([]byte("445918060948226060"), []byte("1234-5678-9123"))
		err = b1.Put([]byte("184576553089368064"), []byte("3141-3285-0154"))
		err = b1.Put([]byte("282501061321818112"), []byte("2540-0360-8159"))
		err = b1.Put([]byte("273627921758027776"), []byte("5117-0178-2284"))
		err = b1.Put([]byte("553781458574114816"), []byte("4043-2766-0969"))
		err = b1.Put([]byte("256510113773387779"), []byte("3742-6311-9777"))
		err = b1.Put([]byte("521771877870600192"), []byte("3871-4778-6691"))
		err = b1.Put([]byte("425071842718515202"), []byte("3957-3730-6484"))
		err = b1.Put([]byte("694517941122629643"), []byte("4429-7799-8199"))
		err = b1.Put([]byte("294632716174098444"), []byte("4687-5239-9514"))
		err = b1.Put([]byte("193730966978691072"), []byte("3613-7802-3642"))
		err = b1.Put([]byte("464765173467447297"), []byte("5417-6609-4455"))
		err = b1.Put([]byte("426077878472540160"), []byte("2067-5941-6303"))
		err = b1.Put([]byte("318836746328735745"), []byte("2411-1823-7939"))
		err = b1.Put([]byte("474223237258149898"), []byte("2497-0851-5644"))
		s.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Your FC has been set to: %s", fc))
		return err
	}); err != nil {
		log.Println(err)
		return
	}
}
