package compose

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func MCPHandler(event *handler.CommandEvent) error {
	err := event.CreateMessage(discord.MessageCreate{
		Content: "TODO",
	})

	return err
}
