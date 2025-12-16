package test

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
)

var Definition = discord.SlashCommandCreate{
	Name:        "componenttest",
	Description: "Command to send kitchen sink component message for feature testing",
}

func Handler(event *handler.CommandEvent) error {
	err := event.CreateMessage(discord.NewMessageCreateBuilder().
		AddActionRow(
			discord.NewButton(discord.ButtonStyleLink, "hey!", "", "https://google.com", snowflake.ID(0)),
		).Build(),
	)

	return err
}
