package commands

import (
	"blast/api"
	"blast/db"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var template = discord.SlashCommandCreate{
	Name:        "template",
	Description: "template command",
}

var Template = Command{
	Handler: func(event *handler.CommandEvent, blast api.EpicClient, user db.UserEntry, credentials api.UserCredentialsResponse, data discord.SlashCommandInteractionData) error {
		time.Sleep(5 * time.Second)

		_, err := event.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
			SetContentf("template").
			AddActionRow(discord.NewPrimaryButton("test", "test_button")).
			Build(),
		)

		return err
	},
	LoginRequired:     true,
	EphemeralResponse: false,
}

// https://user-search-service-prod.ol.epicgames.com/api/v1/search
