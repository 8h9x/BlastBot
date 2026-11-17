package showtoken

import (
	"fmt"
	"time"

	"github.com/8h9x/BlastBot/internal/sessions"
	"github.com/8h9x/fortgo/account"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
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

	exchange, err := session.AccountService.GetExchangeCode()
	if err != nil {
		return err
	}

	token, err := session.AccountService.GetToken(
		&account.AuthClient{
			ID:     "xyza7891p5D7s9R6Gm6moTHWGloerp7B",
			Secret: "Knh18du4NVlFs+3uQ+ZPpDCVto0WYf4yXP8+OcwVt1o",
		},
		account.TokenGrantExchangeCode(exchange.Code),
		account.TokenTypeEG1,
	)

	embed := discord.NewEmbed().
		//        WithAuthorIcon(avatarURL).
		WithColor(0xFB5A32).
		WithTimestamp(time.Now()).
		//        WithAuthorNamef("Launch Fortnite on Windows as %s", me.DisplayName).
		//        WithTitle("Type or paste the following text into a Windows Command Prompt (cmd.exe) and press `Enter`. Expires in 5 minutes.").
		WithDescriptionf("```%s```", token.AccessToken)

	err = event.CreateMessage(discord.NewMessageCreate().WithEphemeral(true).WithEmbeds(embed))
	if err != nil {
		return err
	}

	return nil
}
