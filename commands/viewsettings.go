package commands

import (
	"boobot/structs"
	"boobot/utils"
	"fmt"
	"reflect"

	"github.com/bwmarrin/discordgo"
)

func init() {
	CommandHandlers["viewsettings"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// This is redundant but I'm paranoid
		hasPerm, _ := utils.MemberHasPermission(s, i.GuildID, i.Member.User.ID, discordgo.PermissionAdministrator)
		if !hasPerm {
			respond(s, i, "You don't have the permission to use this command.", true)
			return
		}

		settings := structs.GetSettings(i.GuildID)

		msg := "```asciidoc\n"
		msg += "= Current guild settings ="

		k := reflect.Indirect(reflect.ValueOf(settings))
		typeOfT := k.Type()

		for i := 1; i < k.NumField(); i++ {
			f := k.Field(i)
			msg += fmt.Sprintf("\n%s :: %v", typeOfT.Field(i).Name, f.Interface())
		}
		msg += "\n/changesetting <SettingName> <Value> to change a setting."
		msg += "\nGo to https://boobot.glitch.me/help.html for more info on how to setup these guild settings."
		msg += "```"
		respond(s, i, msg, true)
	}
}
