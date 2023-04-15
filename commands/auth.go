package commands

import (
	"blast/api"
	"blast/api/consts"
	"blast/db"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var auth = discord.SlashCommandCreate{
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
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        "client",
					Description: "The auth client to use.",
					Choices: []discord.ApplicationCommandOptionChoiceString{
						{Name: "fortniteIOSGameClient", Value: generateClientToken(consts.FORTNITE_IOS_CLIENT_ID, consts.FORTNITE_IOS_CLIENT_SECRET)},
						{Name: "launcherAppClient2", Value: generateClientToken(consts.LAUNCHER_CLIENT_ID, consts.LAUNCHER_CLIENT_SECRET)},
						{Name: "wexPCGameClient", Value: generateClientToken(consts.WEX_PC_CLIENT_ID, consts.WEX_PC_CLIENT_SECRET)},
						{Name: "UEFN", Value: generateClientToken(consts.UEFN_CLIENT_ID, consts.UEFN_CLIENT_SECRET)},
						{Name: "fortniteAndroidGameClient", Value: generateClientToken(consts.FORTNITE_ANDROID_CLIENT_ID, consts.FORTNITE_ANDROID_CLIENT_SECRET)},
						{Name: "fortniteSwitchGameClient", Value: generateClientToken(consts.FORTNITE_SWITCH_CLIENT_ID, consts.FORTNITE_SWITCH_CLIENT_SECRET)},
						{Name: "oceanIOSGameClient", Value: generateClientToken(consts.OCEAN_IOS_CLIENT_ID, consts.OCEAN_IOS_CLIENT_SECRET)},
						{Name: "Dauntless", Value: generateClientToken(consts.DAUNTLESS_CLIENT_ID, consts.DAUNTLESS_CLIENT_SECRET)},
						{Name: "wexAndroidGameClient", Value: generateClientToken(consts.WEX_ANDROID_CLIENT_ID, consts.WEX_ANDROID_CLIENT_SECRET)},
						{Name: "fortnitePCGameClient", Value: generateClientToken(consts.FORTNITE_PC_CLIENT_ID, consts.FORTNITE_PC_CLIENT_SECRET)},
						{Name: "wexIOSGameClient", Value: generateClientToken(consts.WEX_IOS_CLIENT_ID, consts.WEX_IOS_CLIENT_SECRET)},
						{Name: "Rumbleverse", Value: generateClientToken(consts.RUMBLEVERSE_CLIENT_ID, consts.RUMBLEVERSE_CLIENT_SECRET)},
						{Name: "Predecessor", Value: generateClientToken(consts.PREDECESSOR_CLIENT_ID, consts.PREDECESSOR_CLIENT_SECRET)},
						{Name: "High On Life", Value: generateClientToken(consts.HIGHONLIFE_CLIENT_ID, consts.HIGHONLIFE_CLIENT_SECRET)},
						{Name: "Bethesda", Value: generateClientToken(consts.BETHESDA_CLIENT_ID, consts.BETHESDA_CLIENT_SECRET)},
						{Name: "Fortnite", Value: generateClientToken(consts.FORTNITE_CLIENT_ID, consts.FORTNITE_CLIENT_SECRET)},
						{Name: "PUBG", Value: generateClientToken(consts.PUBG_CLIENT_ID, consts.PUBG_CLIENT_SECRET)},
						{Name: "MultiVersus", Value: generateClientToken(consts.MULTIVERSUS_CLIENT_ID, consts.MULTIVERSUS_CLIENT_SECRET)},
						{Name: "Twinmotion", Value: generateClientToken(consts.TWINMOTION_CLIENT_ID, consts.TWINMOTION_CLIENT_SECRET)},
						{Name: "Paragon: The Overprime", Value: generateClientToken(consts.PARAGONOVERPRIME_CLIENT_ID, consts.PARAGONOVERPRIME_CLIENT_SECRET)},
						{Name: "Core", Value: generateClientToken(consts.CORE_CLIENT_ID, consts.CORE_CLIENT_SECRET)},
						{Name: "Fall Guys", Value: generateClientToken(consts.FALLGUYS_CLIENT_ID, consts.FALLGUYS_CLIENT_SECRET)},
						{Name: "Rocket League", Value: generateClientToken(consts.ROCKETLEAGUE_CLIENT_ID, consts.ROCKETLEAGUE_CLIENT_SECRET)},
						{Name: "Among Us", Value: generateClientToken(consts.AMONGUS_CLIENT_ID, consts.AMONGUS_CLIENT_SECRET)},
						{Name: "Unreal Editor for Fortnite", Value: generateClientToken(consts.UNREALEDITORFN_CLIENT_ID, consts.UNREALEDITORFN_CLIENT_SECRET)},
					},
					Required: true,
				},
			},
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "device",
			Description: "Generates device authorization using the fortniteIOSGameClient.",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "exchange",
			Description: "Generates an exchange code.",
		},
	},
}

var AuthBearer = Command{
	Handler: func(event *handler.CommandEvent, blast api.EpicClient, user db.UserEntry, credentials api.UserCredentialsResponse, data discord.SlashCommandInteractionData) error {
		_, err := event.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
			SetEmbeds(discord.NewEmbedBuilder().
				SetColor(0xFB5A32).
				SetTimestamp(time.Now()).
				SetTitle("Here's your `fortnitePCGameClient` bearer token:").
				SetDescriptionf("```yml\n%s\n```", credentials.AccessToken).
				Build(),
			).Build(),
		)
		return err
	},
	LoginRequired:     true,
	EphemeralResponse: true,
}

var AuthClient = Command{
	Handler: func(event *handler.CommandEvent, blast api.EpicClient, user db.UserEntry, credentials api.UserCredentialsResponse, data discord.SlashCommandInteractionData) error {
		client := strings.Split(data.String("client"), ":")

		clientCredentials, err := blast.GetClientCredentials(client[0], client[1])
		if err != nil {
			return err
		}

		clientCredentialsJSONStr, err := json.MarshalIndent(clientCredentials, "", "    ")
		if err != nil {
			return err
		}

		_, err = event.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().SetContentf("```json\n%s\n```", clientCredentialsJSONStr).Build())
		return err
	},
	LoginRequired:     false,
	EphemeralResponse: false,
}

var AuthDevice = Command{
	Handler: func(event *handler.CommandEvent, blast api.EpicClient, user db.UserEntry, credentials api.UserCredentialsResponse, data discord.SlashCommandInteractionData) error {
		exchange, err := blast.GetExchangeCode(credentials)
		if err != nil {
			return err
		}

		exchangeCredentials, err := blast.ExchangeCodeLogin(consts.FORTNITE_IOS_CLIENT_ID, consts.FORTNITE_IOS_CLIENT_SECRET, exchange.Code)
		if err != nil {
			return err
		}

		deviceAuth, err := blast.GetDeviceAuth(exchangeCredentials)
		if err != nil {
			return err
		}

		deviceAuthJSONStr, err := json.MarshalIndent(deviceAuth, "", "    ")
		if err != nil {
			return err
		}

		_, err = event.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().SetContentf("```json\n%s\n```", deviceAuthJSONStr).Build())
		return err
	},
	LoginRequired:     true,
	EphemeralResponse: true,
}

var AuthExchange = Command{
	Handler: func(event *handler.CommandEvent, blast api.EpicClient, user db.UserEntry, credentials api.UserCredentialsResponse, data discord.SlashCommandInteractionData) error {
		exchange, err := blast.GetExchangeCode(credentials)
		if err != nil {
			return err
		}

		_, err = event.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
			SetEmbeds(discord.NewEmbedBuilder().
				SetColor(0xFB5A32).
				SetTimestamp(time.Now()).
				SetTitle("Here's your exchange code:").
				SetDescriptionf("```yml\n%s\n```", exchange.Code).
				Build(),
			).Build(),
		)
		return err
	},
	LoginRequired:     true,
	EphemeralResponse: true,
}

func generateClientToken(clientId string, secret string) string {
	return fmt.Sprintf("%s:%s", clientId, secret)
}
