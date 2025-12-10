package commands

import "github.com/disgoorg/disgo/discord"

var logout = discord.SlashCommandCreate{
	Name:        "logout",
	Description: "Removes currently selected account from the bot. (alias for /accounts remove {current})",
}
