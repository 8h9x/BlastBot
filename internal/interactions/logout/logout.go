package logout

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var Definition = discord.SlashCommandCreate{
	Name:        "logout",
	Description: "Removes currently selected account from the bot. (alias for /accounts remove {current})",
}

func Handler(event *handler.CommandEvent) error {
	err := event.CreateMessage(discord.MessageCreate{
		Content: "TODO",
	})

	return err
}
