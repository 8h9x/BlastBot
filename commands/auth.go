package commands

import (
	"blast/api/consts"
	"blast/db"
	"fmt"
	"log"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"go.mongodb.org/mongo-driver/bson"
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
	Handler: func(event *events.ApplicationCommandInteractionCreate) error {
		data := event.SlashCommandInteractionData()

		userId := event.User().ID.String()

		user, err := db.Fetch[db.UserEntry]("users", bson.M{"discordId": userId})
		if err != nil {
			return err
		}

		refreshCredentials, err := blast.RefreshTokenLogin(consts.FORTNITE_PC_CLIENT_ID, consts.FORTNITE_PC_CLIENT_SECRET, user.Accounts[user.SelectedAccount].RefreshToken)
		if err != nil {
			return err
		}

		switch *data.SubCommandName {
		case "bearer": // TODO: make ephermal
			return fmt.Errorf(refreshCredentials.AccessToken)
		case "client":
			clientCredentials, err := blast.GetClientCredentials(consts.FORTNITE_PC_CLIENT_ID, consts.FORTNITE_PC_CLIENT_SECRET)
			if err != nil {
				return err
			}

			return fmt.Errorf(clientCredentials.AccessToken)
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

			log.Println(deviceAuth)

			return fmt.Errorf("done")
		case "exchange": // TODO: make ephermal
			exchange, err := blast.GetExchangeCode(refreshCredentials)
			if err != nil {
				return err
			}

			return fmt.Errorf(exchange.Code)
		default:
			return fmt.Errorf("unknown subcommand")
		}
	},
	LoginRequired: true,
}

// ApplicationCommandOptionSubCommand
