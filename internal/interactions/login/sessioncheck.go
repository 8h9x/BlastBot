package login

import (
	"fmt"

	"github.com/8h9x/BlastBot/database/internal/manager/sessions"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var SessionCheckDefinition = discord.SlashCommandCreate{
	Name:        "sessioncheck",
	Description: "checks current session...",
}

func SessionCheckHandler(event *handler.CommandEvent) error {
	discordId := event.User().ID

	session, err := sessions.GetSessionForUser(discordId)
	if err != nil {
		return fmt.Errorf("unable to get session for user (%s): %s", discordId, err)
	}

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
