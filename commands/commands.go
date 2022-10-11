package commands

import (
	"github.com/bwmarrin/discordgo"
)

var noDM = false
var adminPerms = int64(discordgo.PermissionAdministrator)

var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        "choose",
		Description: "Randomly pick an item from a list containing 2-10 items",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "item-1",
				Description: "Item 1",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "item-2",
				Description: "Item 2",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "item-3",
				Description: "Item 3",
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "item-4",
				Description: "Item 4",
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "item-5",
				Description: "Item 5",
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "item-6",
				Description: "Item 6",
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "item-7",
				Description: "Item 7",
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "item-8",
				Description: "Item 8",
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "item-9",
				Description: "Item 9",
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "item-10",
				Description: "Item 10",
			},
		},
	},
	{
		Name:        "fc",
		Description: "Get the friend code of a user",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "User",
				Required:    true,
			},
		},
		DMPermission: &noDM,
	},
	{
		Name:        "setfc",
		Description: "Set your friend code for this server",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "fc",
				Description: "Your 12 digit friend code",
				Required:    true,
			},
		},
		DMPermission: &noDM,
	},
	{
		Name:         "remfc",
		Description:  "Delete your friend code for this server",
		DMPermission: &noDM,
	},
	{
		Name:        "mmr",
		Description: "Get the MMR/Elo of up to 12 users",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "tracklist",
				Description: "Check the RT or CT leaderboard. Select any option if you don't have multiple leaderboards.",
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "Regular Tracks",
						Value: "rt",
					},
					{
						Name:  "Custom Tracks",
						Value: "ct",
					},
				},
				Required: true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user-1",
				Description: "Name of user 1",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user-2",
				Description: "Name of user 2",
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user-3",
				Description: "Name of user 3",
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user-4",
				Description: "Name of user 4",
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user-5",
				Description: "Name of user 5",
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user-6",
				Description: "Name of user 6",
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user-7",
				Description: "Name of user 7",
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user-8",
				Description: "Name of user 8",
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user-9",
				Description: "Name of user 9",
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user-10",
				Description: "Name of user 10",
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user-11",
				Description: "Name of user 11",
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user-12",
				Description: "Name of user 12",
			},
		},
		DMPermission: &noDM,
	},
	{
		Name:        "stats",
		Description: "Get the full stats of a user",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "tracklist",
				Description: "Check the RT or CT leaderboard. Select any option if you don't have multiple leaderboards.",
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "Regular Tracks",
						Value: "rt",
					},
					{
						Name:  "Custom Tracks",
						Value: "ct",
					},
				},
				Required: true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user-1",
				Description: "Name of user",
				Required:    true,
			},
		},
		DMPermission: &noDM,
	},
	{
		Name:                     "viewsettings",
		Description:              "View all guild settings. Requires Administrator",
		DefaultMemberPermissions: &adminPerms,
		DMPermission:             &noDM,
	},
	{
		Name:        "changesetting",
		Description: "Change a guild setting. Requires Administrator",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "setting",
				Description: "Setting that you want to change.",
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "GameBoards1",
						Value: "GameBoards1",
					},
					{
						Name:  "GameBoards2",
						Value: "GameBoards2",
					},
					{
						Name:  "Spreadsheet1",
						Value: "Spreadsheet1",
					},
					{
						Name:  "Spreadsheet2",
						Value: "Spreadsheet2",
					},
					{
						Name:  "SheetName",
						Value: "SheetName",
					},
					{
						Name:  "RatingName",
						Value: "RatingName",
					},
					{
						Name:  "PlayerIndex",
						Value: "PlayerIndex",
					},
					{
						Name:  "RatingIndex",
						Value: "RatingIndex",
					},
					{
						Name:  "StatIndexes",
						Value: "StatIndexes",
					},
				},
				Required: true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "new-value",
				Description: "New value for the selected setting.",
				Required:    true,
			},
		},
		DefaultMemberPermissions: &adminPerms,
		DMPermission:             &noDM,
	},
	{
		Name:        "clearsetting",
		Description: "Clear the value of a guild setting. Requires Administrator",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "setting",
				Description: "Setting that you want to change.",
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "GameBoards1",
						Value: "GameBoards1",
					},
					{
						Name:  "GameBoards2",
						Value: "GameBoards2",
					},
					{
						Name:  "Spreadsheet1",
						Value: "Spreadsheet1",
					},
					{
						Name:  "Spreadsheet2",
						Value: "Spreadsheet2",
					},
					{
						Name:  "SheetName",
						Value: "SheetName",
					},
					{
						Name:  "RatingName",
						Value: "RatingName",
					},
					{
						Name:  "PlayerIndex",
						Value: "PlayerIndex",
					},
					{
						Name:  "RatingIndex",
						Value: "RatingIndex",
					},
					{
						Name:  "StatIndexes",
						Value: "StatIndexes",
					},
				},
				Required: true,
			},
		},
		DefaultMemberPermissions: &adminPerms,
		DMPermission:             &noDM,
	},
}

// Map of commands that will be added to on startup
var CommandHandlers = make(map[string]func(*discordgo.Session, *discordgo.InteractionCreate))

// Helper function to respond to commands
func respond(s *discordgo.Session, i *discordgo.InteractionCreate, msg string, private bool) {
	if private {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: msg,
			},
		})
	} else {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
			},
		})
	}
}

// Helper function to respond to commands with an embed
func respondEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, embed *discordgo.MessageEmbed, private bool) {
	if private {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:  discordgo.MessageFlagsEphemeral,
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		})
	} else {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		})
	}
}
