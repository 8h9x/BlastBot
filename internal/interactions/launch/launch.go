package launch

import (
	"fmt"
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
	err := event.DeferCreateMessage(true)
	if err != nil {
		return err
	}

	discordId := event.User().ID

	session, err := sessions.GetSessionForUser(discordId)
	if err != nil {
		return fmt.Errorf("unable to get session for user (%s): %s", discordId, err)
	}

	me, err := session.FetchMe()
	if err != nil {
		return fmt.Errorf("unable to fetch me: %s", err)
	}

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

	flags := discord.MessageFlagIsComponentsV2
	_, err = event.UpdateInteractionResponse(discord.MessageUpdate{
		Flags: &flags,
		Components: &[]discord.LayoutComponent{
			discord.NewContainer(
				discord.NewTextDisplayf("## Launch Fortnite on Windows as %s", me.DisplayName),
				discord.NewTextDisplay("Paste the following text into a Windows Command Prompt (cmd.exe) and press `Enter`."),
				discord.NewTextDisplayf("```powershell\nstart /d \"%%PROGRAMFILES%%\\Epic Games\\Fortnite\\FortniteGame\\Binaries\\Win64\" FortniteLauncher.exe -AUTH_LOGIN=unused -AUTH_PASSWORD=%s -AUTH_TYPE=exchangecode -epicapp=Fortnite -epicenv=Prod  -EpicPortal -epicsandboxid=fn -epicuserid=%s\n```", launcherExchange.Code, exchangeCredentials.AccountID),
				discord.NewTextDisplayf("Expires <t:%d:R>", time.Now().Add(5 * time.Minute).Unix()),
			),
		},
	});
	return err
}