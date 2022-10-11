package commands

import (
	"boobot/structs"
	"boobot/utils"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func init() {
	CommandHandlers["changesetting"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// This is redundant but I'm paranoid
		hasPerm, _ := utils.MemberHasPermission(s, i.GuildID, i.Member.User.ID, discordgo.PermissionAdministrator)
		if !hasPerm {
			respond(s, i, "You don't have the permission to use this command.", true)
			return
		}

		settings := structs.GetSettings(i.GuildID)

		setting := i.ApplicationCommandData().Options[0].StringValue()
		newval := i.ApplicationCommandData().Options[1].StringValue()
		valid := settings.Set(setting, newval)
		if !valid {
			// Should never happen
			respond(s, i, fmt.Sprintf("`%s` is not a valid setting.", setting), true)
			return
		}
		respond(s, i, fmt.Sprintf("Value of setting `%s` changed to: `%s`", setting, newval), true)
	}

	CommandHandlers["clearsetting"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// This is redundant but I'm paranoid
		hasPerm, _ := utils.MemberHasPermission(s, i.GuildID, i.Member.User.ID, discordgo.PermissionAdministrator)
		if !hasPerm {
			respond(s, i, "You don't have the permission to use this command.", true)
			return
		}

		settings := structs.GetSettings(i.GuildID)

		setting := i.ApplicationCommandData().Options[0].StringValue()
		valid := settings.Set(setting, "")
		if !valid {
			// Should never happen
			respond(s, i, fmt.Sprintf("`%s` is not a valid setting.", setting), true)
			return
		}
		respond(s, i, fmt.Sprintf("Value of setting `%s` has been cleared.", setting), true)
	}
}
