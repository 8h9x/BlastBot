package item_shop

import "github.com/disgoorg/disgo/discord"

var Definition = discord.SlashCommandCreate{
	Name:        "item-shop",
	Description: "Collection of commands to manage linked accounts.",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionSubCommand{
			Name:        "activity",
			Description: "Launches a discord activity for the current item shop.",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "buy",
			Description: "Quickly purchase an item from the item shop.",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "image",
			Description: "Generates an image of the current item shop.",
		},
	},
}