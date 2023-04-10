package commands

import (
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var auth = Command{
	Create: discord.SlashCommandCreate{
		Name:        "auth",
		Description: "Authentication related commands.",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionSubCommand{
				Name:        "bearer",
				Description: "Generates a bearer token.",
			},
			discord.ApplicationCommandOptionSubCommand{
				Name:        "client",
				Description: "Generates client credentials.",
			},
			discord.ApplicationCommandOptionSubCommand{
				Name:        "device",
				Description: "Generates device authorization.",
			},
			discord.ApplicationCommandOptionSubCommand{
				Name:        "exchange",
				Description: "Generates an exchange code.",
			},
			discord.ApplicationCommandOptionSubCommand{
				Name:        "refresh",
				Description: "Generates a refresh token",
			},
		},
	},
	Handler: func(event *events.ApplicationCommandInteractionCreate) error {
		embed := discord.NewEmbedBuilder().
			// SetAuthorIcon(event.). // TODO set author icon to bot user avatar
			SetColor(0xFB5A32).
			SetTimestamp(time.Now()).
			SetTitle("Test").
			SetDescription("test").
			Build()

		err := event.CreateMessage(discord.NewMessageCreateBuilder().
			SetEmbeds(embed).
			AddActionRow(
				discord.NewDangerButton("test", "cancel-login"),
			).
			Build(),
		)
		if err != nil {
			return err
		}

		return nil
	},
}

// ApplicationCommandOptionSubCommand
