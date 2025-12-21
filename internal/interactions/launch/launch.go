package launch

import (
	"fmt"
	"strings"
	"time"

    "github.com/8h9x/BlastBot/internal/sessions"
	"github.com/8h9x/fortgo/auth"
	"github.com/8h9x/fortgo/consts"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var Definition = discord.SlashCommandCreate{
	Name:        "launch",
	Description: "Produces a windows cmd string to launch fortnite using the current account.",
}

func Handler(event *handler.CommandEvent) error {
	discordId := event.User().ID

	session, err := sessions.GetSessionForUser(discordId)
	if err != nil {
		return fmt.Errorf("unable to get session for user (%s): %s", discordId, err)
	}

	me, err := session.FetchMe()
	if err != nil {
		return fmt.Errorf("unable to fetch me: %s", err)
	}

	avatar, err := session.AvatarService.GetOne(me.ID)
    if err != nil {
        return err
    }

	avatarURL := fmt.Sprintf("https://fortnite-api.com/images/cosmetics/br/%s/icon.png", strings.Replace(avatar.AvatarID, "ATHENACHARACTER:", "", -1))

	exchange, err := session.GetExchangeCode()
    if err != nil {
        return err
    }

	exchangeCredentials, err := auth.Authenticate(session.HTTPClient, &auth.AuthClient{consts.LauncherClientID, consts.LauncherClientSecret}, auth.PayloadExchangeCode{exchange.Code}, true)//client.Epic.ExchangeCodeLogin(consts.LAUNCHER_CLIENT_ID, consts.LAUNCHER_CLIENT_SECRET, exchange.Code)
    if err != nil {
        return err
    }

	launcherExchange, err := auth.GetExchangeCode(session.HTTPClient, exchangeCredentials)
	if err != nil {
        return err
    }

    embed := discord.NewEmbedBuilder().
        SetAuthorIcon(avatarURL).
        SetColor(0xFB5A32).
        SetTimestamp(time.Now()).
        SetAuthorNamef("Launch Fortnite on Windows as %s", me.DisplayName).
        SetTitle("Type or paste the following text into a Windows Command Prompt (cmd.exe) and press `Enter`. Expires in 5 minutes.").
        SetDescriptionf("```start /d \"%%PROGRAMFILES%%\\Epic Games\\Fortnite\\FortniteGame\\Binaries\\Win64\" FortniteLauncher.exe -AUTH_LOGIN=unused -AUTH_PASSWORD=%s -AUTH_TYPE=exchangecode -epicapp=Fortnite -epicenv=Prod  -EpicPortal -epicsandboxid=fn -epicuserid=%s```", launcherExchange.Code, exchangeCredentials.AccountID).
        Build()

	err = event.CreateMessage(discord.NewMessageCreateBuilder().SetEphemeral(true).SetEmbeds(embed).Build())
	if err != nil {
		return err
	}

	return nil
}
