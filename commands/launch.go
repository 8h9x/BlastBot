package commands

import (
	"blast/api"
	"blast/db"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var launch = discord.SlashCommandCreate{
	Name:        "launch",
	Description: "Windows command to launch fortnite using the currently selected account.",
}

var Launch = Command{
	Handler: func(event *handler.CommandEvent, blast api.EpicClient, user db.UserEntry, credentials api.UserCredentialsResponse, data discord.SlashCommandInteractionData) error {
		account, err := blast.FetchAccountInformation(credentials)
		if err != nil {
			return err
		}

		avatarURL, err := blast.FetchAvatarURL(credentials)
		if err != nil {
			return err
		}

		exchange, err := blast.GetExchangeCode(credentials)
		if err != nil {
			return err
		}

		embed := discord.NewEmbedBuilder().
			SetAuthorIcon(avatarURL).
			SetColor(0xFB5A32).
			SetTimestamp(time.Now()).
			SetAuthorNamef("Launch Fortnite on Windows as %s", account.DisplayName).
			SetTitle("Type or paste the following text into a Windows Command Prompt (cmd.exe) and press `Enter`. Expires in 5 minutes.").
			SetDescriptionf("```start /d \"C:\\Program Files\\Epic Games\\Fortnite\\FortniteGame\\Binaries\\Win64\" FortniteLauncher.exe -AUTH_LOGIN=unused -AUTH_PASSWORD=%s -AUTH_TYPE=exchangecode -epicapp=Fortnite -epicenv=Prod  -EpicPortal -epicsandboxid=fn -epicuserid=%s```", exchange.Code, credentials.AccountID).
			Build()

		_, err = event.Client().Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.NewMessageUpdateBuilder().SetEmbeds(embed).Build())
		return err
	},
	LoginRequired:     true,
	EphemeralResponse: false,
}

// https://user-search-service-prod.ol.epicgames.com/api/v1/search
