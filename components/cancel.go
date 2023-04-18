package components

import (
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

// TODO: delete message instead
func Cancel(e *handler.ComponentEvent) error {
	return e.UpdateMessage(discord.NewMessageUpdateBuilder().
		SetEmbeds(discord.NewEmbedBuilder().
			SetColor(0xCBA6F7).
			SetTimestamp(time.Now()).
			SetTitle("Canceled!").
			Build()).
		ClearContainerComponents().
		Build(),
	)
}
