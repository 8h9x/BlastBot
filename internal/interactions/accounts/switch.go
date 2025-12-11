package accounts

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func SwitchHandler(event *handler.CommandEvent) error {
	err := event.CreateMessage(discord.MessageCreate{
		Content: "TODO",
	})

	return err
}
