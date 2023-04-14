package commands

import (
	"blast/api/consts"
	"blast/db"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var auth = Command{
	Create: discord.SlashCommandCreate{
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
							{Name: "fortniteIOSGameClient", Value: fmt.Sprintf("%s:%s", consts.FORTNITE_IOS_CLIENT_ID, consts.FORTNITE_IOS_CLIENT_SECRET)},
							{Name: "launcherAppClient2", Value: fmt.Sprintf("%s:%s", consts.LAUNCHER_CLIENT_ID, consts.LAUNCHER_CLIENT_SECRET)},
							{Name: "wexPCGameClient", Value: fmt.Sprintf("%s:%s", consts.WEX_PC_CLIENT_ID, consts.WEX_PC_CLIENT_SECRET)},
							{Name: "UEFN", Value: fmt.Sprintf("%s:%s", consts.UEFN_CLIENT_ID, consts.UEFN_CLIENT_SECRET)},
							{Name: "fortniteAndroidGameClient", Value: fmt.Sprintf("%s:%s", consts.FORTNITE_ANDROID_CLIENT_ID, consts.FORTNITE_ANDROID_CLIENT_SECRET)},
							{Name: "fortniteSwitchGameClient", Value: fmt.Sprintf("%s:%s", consts.FORTNITE_SWITCH_CLIENT_ID, consts.FORTNITE_SWITCH_CLIENT_SECRET)},
							{Name: "oceanIOSGameClient", Value: fmt.Sprintf("%s:%s", consts.OCEAN_IOS_CLIENT_ID, consts.OCEAN_IOS_CLIENT_SECRET)},
							{Name: "Dauntless", Value: fmt.Sprintf("%s:%s", consts.DAUNTLESS_CLIENT_ID, consts.DAUNTLESS_CLIENT_SECRET)},
							{Name: "wexAndroidGameClient", Value: fmt.Sprintf("%s:%s", consts.WEX_ANDROID_CLIENT_ID, consts.WEX_ANDROID_CLIENT_SECRET)},
							{Name: "fortnitePCGameClient", Value: fmt.Sprintf("%s:%s", consts.FORTNITE_PC_CLIENT_ID, consts.FORTNITE_PC_CLIENT_SECRET)},
							{Name: "wexIOSGameClient", Value: fmt.Sprintf("%s:%s", consts.WEX_IOS_CLIENT_ID, consts.WEX_IOS_CLIENT_SECRET)},
							{Name: "Rumbleverse", Value: fmt.Sprintf("%s:%s", consts.RUMBLEVERSE_CLIENT_ID, consts.RUMBLEVERSE_CLIENT_SECRET)},
							{Name: "Predecessor", Value: fmt.Sprintf("%s:%s", consts.PREDECESSOR_CLIENT_ID, consts.PREDECESSOR_CLIENT_SECRET)},
							{Name: "High On Life", Value: fmt.Sprintf("%s:%s", consts.HIGHONLIFE_CLIENT_ID, consts.HIGHONLIFE_CLIENT_SECRET)},
							{Name: "Bethesda", Value: fmt.Sprintf("%s:%s", consts.BETHESDA_CLIENT_ID, consts.BETHESDA_CLIENT_SECRET)},
							{Name: "Fortnite", Value: fmt.Sprintf("%s:%s", consts.FORTNITE_CLIENT_ID, consts.FORTNITE_CLIENT_SECRET)},
							{Name: "PUBG", Value: fmt.Sprintf("%s:%s", consts.PUBG_CLIENT_ID, consts.PUBG_CLIENT_SECRET)},
							{Name: "MultiVersus", Value: fmt.Sprintf("%s:%s", consts.MULTIVERSUS_CLIENT_ID, consts.MULTIVERSUS_CLIENT_SECRET)},
							{Name: "Twinmotion", Value: fmt.Sprintf("%s:%s", consts.TWINMOTION_CLIENT_ID, consts.TWINMOTION_CLIENT_SECRET)},
							{Name: "Paragon: The Overprime", Value: fmt.Sprintf("%s:%s", consts.PARAGONOVERPRIME_CLIENT_ID, consts.PARAGONOVERPRIME_CLIENT_SECRET)},
							{Name: "Core", Value: fmt.Sprintf("%s:%s", consts.CORE_CLIENT_ID, consts.CORE_CLIENT_SECRET)},
							{Name: "Fall Guys", Value: fmt.Sprintf("%s:%s", consts.FALLGUYS_CLIENT_ID, consts.FALLGUYS_CLIENT_SECRET)},
							{Name: "Rocket League", Value: fmt.Sprintf("%s:%s", consts.ROCKETLEAGUE_CLIENT_ID, consts.ROCKETLEAGUE_CLIENT_SECRET)},
							{Name: "Among Us", Value: fmt.Sprintf("%s:%s", consts.AMONGUS_CLIENT_ID, consts.AMONGUS_CLIENT_SECRET)},
							{Name: "Unreal Editor for Fortnite", Value: fmt.Sprintf("%s:%s", consts.UNREALEDITORFN_CLIENT_ID, consts.UNREALEDITORFN_CLIENT_SECRET)},
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
	},
	Handler: func(event *events.ApplicationCommandInteractionCreate, user db.UserEntry) error {
		data := event.SlashCommandInteractionData()

		refreshCredentials, err := blast.RefreshTokenLogin(consts.FORTNITE_PC_CLIENT_ID, consts.FORTNITE_PC_CLIENT_SECRET, user.Accounts[user.SelectedAccount].RefreshToken)
		if err != nil {
			return err
		}

		switch *data.SubCommandName {
		case "bearer": // TODO: make ephermal
			return fmt.Errorf(refreshCredentials.AccessToken)
		case "client":
			client := strings.Split(data.String("client"), ":")

			clientCredentials, err := blast.GetClientCredentials(client[0], client[1])
			if err != nil {
				return err
			}

			clientCredentialsJSONStr, err := json.MarshalIndent(clientCredentials, "", "    ")
			if err != nil {
				return err
			}

			_, err = event.Client().Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.NewMessageUpdateBuilder().SetContentf("```json\n%s\n```", clientCredentialsJSONStr).Build())
			if err != nil {
				return err
			}
		case "device": // TODO: make ephermal
			exchange, err := blast.GetExchangeCode(refreshCredentials)
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

			_, err = event.Client().Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.NewMessageUpdateBuilder().SetContentf("```json\n%s\n```", deviceAuthJSONStr).Build())
			if err != nil {
				return err
			}
		case "exchange": // TODO: make ephermal
			exchange, err := blast.GetExchangeCode(refreshCredentials)
			if err != nil {
				return err
			}

			_, err = event.Client().Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.NewMessageUpdateBuilder().
				SetEmbeds(discord.NewEmbedBuilder().
					SetColor(0xFB5A32).
					SetTimestamp(time.Now()).
					SetTitle("Here's your exchange code:").
					SetDescriptionf("```hs\n%s\n```", exchange.Code).
					Build(),
				).Build(),
			)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown subcommand")
		}

		return nil
	},
	LoginRequired: true,
}

// ApplicationCommandOptionSubCommand
