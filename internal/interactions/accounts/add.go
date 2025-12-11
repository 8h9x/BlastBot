package accounts

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func AddHandler(event *handler.CommandEvent) error {
	err := event.CreateMessage(discord.MessageCreate{
		Content: "TODO",
	})

	return err
}
