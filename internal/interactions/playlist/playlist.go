package mnemonic

import "github.com/disgoorg/disgo/discord"

var opts = []discord.ApplicationCommandOption{
	discord.ApplicationCommandOptionString{
		Name:        "playlist",
		Description: "Playlist id or creative map code.",
		Required:    true,
	},
}

var Definition = discord.SlashCommandCreate{
	Name:        "playlist",
	Description: "Mnemonic related commands.",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionSubCommandGroup{
			Name:        "favorites",
			Description: "Mnemonic favorites related commands.",
			Options: []discord.ApplicationCommandOptionSubCommand{
				{
					Name:        "add",
					Description: "Adds a playlist to your favorites list.",
					Options:     opts,
				},
				{
					Name:        "list",
					Description: "Returns your discovery favorites list.",
				},
				{
					Name:        "remove",
					Description: "Removes a playlist from your favorites list.",
					Options:     opts,
				},
			},
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "info",
			Description: "Returns information about a playlist.",
			Options:     opts,
		},
	},
}
