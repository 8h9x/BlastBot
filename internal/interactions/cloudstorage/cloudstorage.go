package cloudstorage

import "github.com/disgoorg/disgo/discord"

var Definition = discord.SlashCommandCreate{
	Name:        "cloudstorage",
	Description: "Collection of commands to manage cloudstorage files.",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionSubCommand{
			Name:        "dump",
			Description: "Export selected .sav file.",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "edit",
			Description: "Helper command to modify client settings for selected platform.",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "transfer",
			Description: "Copy .sav file from one platform to another",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "upload",
			Description: "Replace cloudstorage file on the server with uploaded one.",
		},
	},
}
