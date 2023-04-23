package commands

import (
	"blast/api"
	"blast/db"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/rest"
)

var ping = discord.SlashCommandCreate{
	Name:        "ping",
	Description: "Check bot latency.",
}

var Ping = Command{
	Handler: func(event *handler.CommandEvent, blast api.EpicClient, user db.UserEntry, credentials api.UserCredentialsResponse, data discord.SlashCommandInteractionData) error {
		var gatewayPing string
		if event.Client().HasGateway() {
			gatewayPing = event.Client().Gateway().Latency().String()
		}

		eb := discord.NewEmbedBuilder().
			SetTitle("Ping").
			AddField("Rest", "loading...", false).
			AddField("Gateway", gatewayPing, false).
			SetColor(0xFB5A32).
			SetTimestamp(time.Now())

		defer func() {
			var start int64
			_, _ = event.Client().Rest().GetBotApplicationInfo(func(config *rest.RequestConfig) {
				start = time.Now().UnixNano()
			})
			duration := time.Now().UnixNano() - start
			eb.SetField(0, "Rest", time.Duration(duration).String(), false)
			if _, err := event.Client().Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.MessageUpdate{Embeds: &[]discord.Embed{eb.Build()}}); err != nil {
				event.Client().Logger().Error("Failed to update ping embed: ", err)
			}
		}()

		_, err := event.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
			SetEmbeds(eb.Build()).
			Build(),
		)

		return err
	},
	LoginRequired:     true,
	EphemeralResponse: false,
}
