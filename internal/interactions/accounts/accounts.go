package accounts

import "github.com/disgoorg/disgo/discord"

var Definition = discord.SlashCommandCreate{
	Name:        "accounts",
	Description: "Collection of commands to manage linked accounts.",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionSubCommand{
			Name:        "add",
			Description: "Login and add a new account to the bot.",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "switch",
			Description: "Switch currently active account.",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "remove",
			Description: "Remove an account from the bot. (Deletes all stored data for this account)",
		},
	},
}
