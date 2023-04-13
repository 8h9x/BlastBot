package commands

import (
	"fmt"

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
			},
			discord.ApplicationCommandOptionSubCommand{
				Name:        "device",
				Description: "Generates device authorization.",
			},
			discord.ApplicationCommandOptionSubCommand{
				Name:        "exchange",
				Description: "Generates an exchange code.",
			},
		},
	},
	Handler: func(event *events.ApplicationCommandInteractionCreate) error {
		data := event.SlashCommandInteractionData()

		// userId := event.User().ID.String()

		// user, err := db.Fetch[db.UserEntry]("users", bson.M{"discordId": userId})
		// if err != nil {
		// 	return err
		// }

		// refreshCredentials, err := blast.RefreshTokenLogin(consts.FORTNITE_PC_CLIENT_ID, consts.FORTNITE_PC_CLIENT_SECRET, user.Accounts[user.SelectedAccount].RefreshToken)
		// if err != nil {
		// 	return err
		// }

		// log.Println(refreshCredentials)

		switch *data.SubCommandName {
		case "bearer":

		case "client":

		case "device":

		case "exchange":

		default:
			return fmt.Errorf("unknown subcommand")
		}

		return nil
	},
	LoginRequired: true,
}

// ApplicationCommandOptionSubCommand
