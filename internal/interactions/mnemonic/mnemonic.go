package mnemonic

import "github.com/disgoorg/disgo/discord"

var opts = []discord.ApplicationCommandOption{
	discord.ApplicationCommandOptionString{
		Name:        "mnemonic",
		Description: "Playlist id or creative map code.",
		Required:    true,
	},
}

var Definition = discord.SlashCommandCreate{
	Name:        "mnemonic",
	Description: "Mnemonic related commands.",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionSubCommandGroup{
			Name:        "favorites",
			Description: "Mnemonic favorites related commands.",
			Options: []discord.ApplicationCommandOptionSubCommand{
				{
					Name:        "add",
					Description: "Adds a mnemonic to your favorites list.",
					Options:     opts,
				},
				{
					Name:        "list",
					Description: "Returns your discovery favorites list.",
				},
				{
					Name:        "remove",
					Description: "Removes a mnemonic from your favorites list.",
					Options:     opts,
				},
			},
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "info",
			Description: "Returns information about a mnemonic.",
			Options:     opts,
		},
	},
}
