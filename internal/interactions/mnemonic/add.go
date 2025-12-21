package mnemonic

import (
	"encoding/json"
	"fmt"
	"time"

    "github.com/8h9x/BlastBot/internal/sessions"
	"github.com/8h9x/fortgo/links"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)


func FavoriteAddHandler(event *handler.CommandEvent) error {
	discordId := event.User().ID
	data := event.SlashCommandInteractionData()
	mnemonic := data.String("mnemonic")

	session, err := sessions.GetSessionForUser(discordId)
	if err != nil {
		return fmt.Errorf("unable to get session for user (%s): %s", discordId, err)
	}

	mnemonicData, err := session.GetMnemonic(mnemonic, links.MnemonicType(""), -1)
	if err != nil {
		return fmt.Errorf("unable to fetch mnemonic data: %s", err)
	}

	mnemonicDataRaw, err := json.Marshal(mnemonicData)
	if err != nil {
		return fmt.Errorf("unable to marshal mnemonic data: %s", err)
	}

	embed := discord.NewEmbedBuilder().
		SetColor(0xFB2C36).
		SetTimestamp(time.Now()).
		SetTitle("<:exclamation:1096641657396539454> We hit a roadblock!").
		SetDescriptionf("```json\n%s\n```", mnemonicDataRaw)

	event.CreateMessage(discord.NewMessageCreateBuilder().
		SetEmbeds(embed.
			Build(),
		).
		Build(),
	)

	me, err := session.FetchMe()
	if err != nil {
		return fmt.Errorf("unable to fetch me: %s", err)
	}

	err = event.CreateMessage(discord.MessageCreate{
		Content: fmt.Sprintf("session using client: %s exists for account: %s, discord user: %s", session.ClientID, me.ID, event.User().ID.String()),
	})
	if err != nil {
		return err
	}

	return nil
}