package utils

import (
	"bytes"
	"fmt"
	"math/rand"
	"reflect"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Prevent bot from crashing from a panic
func RecoverPanic(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "The command crashed while running. Please try again in a moment.",
		},
	})
	if r := recover(); r != nil {
		fmt.Println("RECOVERED PANIC:", r)
		dmChannel, err := s.UserChannelCreate("397514708736802816")
		if err != nil {
			fmt.Println("Error sending panic stack trace DM:\n", err)
		} else {
			msg := fmt.Sprintf("RECOVERED PANIC:\n```%v\n%s```", r, string(debug.Stack()))
			s.ChannelMessageSend(dmChannel.ID, msg)
		}
	}
}

// Contains checks if the slice/array contains the value
func Contains(slice, val interface{}) bool {
	sliceVal := reflect.ValueOf(slice)
	for i := 0; i < sliceVal.Len(); i++ {
		if val == sliceVal.Index(i).Interface() {
			return true
		}
	}
	return false
}

// Select random value from slice/array
func RandomVal(slice interface{}) interface{} {
	slc := reflect.ValueOf(slice)
	// Seed the generator using the current time
	rand.Seed(time.Now().UnixNano())
	val := slc.Index(rand.Intn(slc.Len())).Interface()
	return val
}

// Inserts a string every n characters
func InsertNth(s string, n int) string {
	var buffer bytes.Buffer
	var n_1 = n - 1
	var l_1 = len(s) - 1
	for i, rune := range s {
		buffer.WriteRune(rune)
		if i%n == n_1 && i != l_1 {
			buffer.WriteRune('-')
		}
	}
	return buffer.String()
}

func Ternary(statement bool, a, b interface{}) interface{} {
	if statement {
		return a
	}
	return b
}

func Nth(i int) string {
	if i%10 == 1 && i != 11 {
		return strconv.Itoa(i) + "st"
	}
	if i%10 == 2 && i != 12 {
		return strconv.Itoa(i) + "nd"
	}
	if i%10 == 3 && i != 13 {
		return strconv.Itoa(i) + "rd"
	}
	return strconv.Itoa(i) + "th"
}

func GetArgs(options []*discordgo.ApplicationCommandInteractionDataOption) []string {
	args := []string{}
	for _, option := range options {
		args = append(args, option.StringValue())
	}
	return args
}

func MemberHasPermission(s *discordgo.Session, guildID string, userID string, permission int64) (bool, error) {
	// Guild owner has every permission
	guild, err := s.Guild(guildID)
	if err != nil {
		if guild.OwnerID == userID {
			return true, nil
		}
	}

	member, err := s.State.Member(guildID, userID)
	if err != nil {
		if member, err = s.GuildMember(guildID, userID); err != nil {
			return false, err
		}
	}

	// Iterate through the role IDs stored in member.Roles
	// to check permissions
	for _, roleID := range member.Roles {
		role, err := s.State.Role(guildID, roleID)
		if err != nil {
			return false, err
		}
		if role.Permissions&permission != 0 {
			return true, nil
		}
	}

	return false, nil
}
