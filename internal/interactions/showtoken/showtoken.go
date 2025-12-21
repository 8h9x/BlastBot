package showtoken

import (
    "fmt"
    "time"

    "github.com/8h9x/BlastBot/internal/sessions"
    "github.com/8h9x/fortgo/auth"
    "github.com/disgoorg/disgo/discord"
    "github.com/disgoorg/disgo/handler"
    //	"github.com/8h9x/fortgo/consts"
)

var Definition = discord.SlashCommandCreate{
	Name:        "showtoken",
	Description: "...",
}

func Handler(event *handler.CommandEvent) error {
	discordId := event.User().ID

	session, err := sessions.GetSessionForUser(discordId)
	if err != nil {
		return fmt.Errorf("unable to get session for user (%s): %s", discordId, err)
	}

	exchange, err := session.GetExchangeCode()
    if err != nil {
        return err
    }

	exchangeCredentials, err := auth.Authenticate(session.HTTPClient, &auth.AuthClient{"xyza7891p5D7s9R6Gm6moTHWGloerp7B", "Knh18du4NVlFs+3uQ+ZPpDCVto0WYf4yXP8+OcwVt1o"}, auth.PayloadExchangeCode{exchange.Code}, false)//client.Epic.ExchangeCodeLogin(consts.LAUNCHER_CLIENT_ID, consts.LAUNCHER_CLIENT_SECRET, exchange.Code)
    if err != nil {
        return err
    }

    embed := discord.NewEmbedBuilder().
//        SetAuthorIcon(avatarURL).
        SetColor(0xFB5A32).
        SetTimestamp(time.Now()).
//        SetAuthorNamef("Launch Fortnite on Windows as %s", me.DisplayName).
//        SetTitle("Type or paste the following text into a Windows Command Prompt (cmd.exe) and press `Enter`. Expires in 5 minutes.").
        SetDescriptionf("```%s```", exchangeCredentials.AccessToken).
        Build()

	err = event.CreateMessage(discord.NewMessageCreateBuilder().SetEphemeral(true).SetEmbeds(embed).Build())
	if err != nil {
		return err
	}

	return nil
}
