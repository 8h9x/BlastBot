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

var opts = []discord.ApplicationCommandOption{
	discord.ApplicationCommandOptionString{
		Name:        "mnemonic",
		Description: "Playlist id or creative map code.",
		Required:    true,
	},
}

var mnemonic = Command{
	Create: discord.SlashCommandCreate{
		Name:        "mnemonic",
		Description: "Mnemonic related commands.",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionSubCommand{
				Name:        "favorite",
				Description: "Adds a mnemonic to your favorites.",
				Options:     opts,
			},
			discord.ApplicationCommandOptionSubCommand{
				Name:        "info",
				Description: "Returns information about a mnemonic.",
				Options:     opts,
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

		log.Println(refreshCredentials)

		switch *data.SubCommandName {
		case "favorite":
			err := blast.AddFavoriteMnemonic(refreshCredentials, data.String("mnemonic"))
			if err != nil {
				return err
			}

			return fmt.Errorf("added favorite")
		case "info":
			res, err := blast.FetchMnemonicInfo(refreshCredentials, data.String("mnemonic"))
			if err != nil {
				return err
			}

			log.Println(res)

			return fmt.Errorf("data logged to console")
		default:
			return fmt.Errorf("unknown subcommand")
		}
	},
	LoginRequired: true,
}

// /mnemonic favorites list
// /mnemonic favorites add
// /mnemonic favorites remove
