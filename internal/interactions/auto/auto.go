package auto

import (
	"github.com/disgoorg/disgo/discord"
)

var sharedOpts = []discord.ApplicationCommandOption{
	discord.ApplicationCommandOptionBool{
		Name:        "bulk",
		Description: "Applies the auto setting to all linked accounts.",
		Required:    false,
	},
}

var Definition = discord.SlashCommandCreate{
	Name:        "auto",
	Description: "Collection of commands to enable automatic daily actions/alerts.",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionSubCommand{
			Name:        "winterfest",
			Description: "Automatically open winterfest gifts each day.",
			Options: sharedOpts,
		},
	},
}