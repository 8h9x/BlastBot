package compose

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func RawHandler(event *handler.CommandEvent) error {
	err := event.CreateMessage(discord.MessageCreate{
		Content: "TODO",
	})

	return err
}
