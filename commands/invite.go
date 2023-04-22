package commands

import (
	"blast/api"
	"blast/db"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var invite = discord.SlashCommandCreate{
	Name:        "invite",
	Description: "Sends a link to invite Blast! to your own server.",
}

var Invite = Command{
	Handler: func(event *handler.CommandEvent, blast api.EpicClient, user db.UserEntry, credentials api.UserCredentialsResponse, data discord.SlashCommandInteractionData) error {
		_, err := event.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
			AddActionRow(discord.NewLinkButton("Invite Blast!", "https://discord.com/api/oauth2/authorize?client_id=1099142031680483442&permissions=532576267328&scope=bot%20applications.commands"), discord.NewLinkButton("Support Server", "https://discord.gg/astra-921104988363694130")).
			Build(),
		)

		return err
	},
	LoginRequired:     false,
	EphemeralResponse: false,
}

// https://user-search-service-prod.ol.epicgames.com/api/v1/search
