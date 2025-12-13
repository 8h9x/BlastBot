package compose

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func Definition() discord.ApplicationCommandCreate {
	return discord.SlashCommandCreate{
		Name:        "compose",
		Description: "Advanced command for power users to have more direct access to API calls.",
	}
}

func Handler(event *handler.CommandEvent) error {
	err := event.CreateMessage(discord.MessageCreate{
		Content: "TODO",
	})

	return err
}
