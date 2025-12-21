package redeem

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/8h9x/BlastBot/internal/sessions"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var Definition = discord.SlashCommandCreate{
	Name:        "redeem",
	Description: "Redeem a code.",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:        "code",
			Description: "The code to redeem.",
			Required:    true,
		},
	},
}

func Handler(event *handler.CommandEvent) error {
	discordId := event.User().ID
	data := event.SlashCommandInteractionData()
	code := strings.ToLower(data.String("code"))

	session, err := sessions.GetSessionForUser(discordId)
	if err != nil {
		return fmt.Errorf("unable to get session for user (%s): %s", discordId, err)
	}

	codeData, err := session.FulfillmentService.RedeemCodeForAccount(session.CurrentCredentials().AccountID, code)
	if err != nil {
		return fmt.Errorf("unable to redeem code: %s", err)
	}

	codeDataRaw, err := json.Marshal(codeData)
	if err != nil {
		return fmt.Errorf("unable to marshal code data: %s", err)
	}

	embed := discord.NewEmbedBuilder().
		SetColor(0xFFFFFF).
		SetTimestamp(time.Now()).
		SetTitle("guh!").
		SetDescriptionf("```json\n%s\n```", codeDataRaw)

	err = event.CreateMessage(discord.NewMessageCreateBuilder().
		SetEmbeds(embed.
			Build(),
		).
		Build(),
	)
	if err != nil {
		return err
	}

	return nil
}