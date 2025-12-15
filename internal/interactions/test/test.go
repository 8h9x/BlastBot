package test

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var Definition = discord.SlashCommandCreate{
	Name:        "componenttest",
	Description: "Command to send kitchen sink component message for feature testing",
}

func Handler(event *handler.CommandEvent) error {
	err := event.CreateMessage(discord.NewMessageCreateBuilder().
		AddActionRow(

		)
	)

	return err
}
