package commands

import (
	"blast/api/consts"
	"blast/db"
	"fmt"
	"log"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
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
			discord.ApplicationCommandOptionSubCommandGroup{
				Name:        "favorites",
				Description: "Mnemonic favorites related commands.",
				Options: []discord.ApplicationCommandOptionSubCommand{
					{
						Name:        "list",
						Description: "Returns your discovery favorites list.",
					},
					{
						Name:        "add",
						Description: "Adds a mnemonic to your favorites list.",
						Options:     opts,
					},
					{
						Name:        "remove",
						Description: "Removes a mnemonic from your favorites list.",
						Options:     opts,
					},
				},
			},
			discord.ApplicationCommandOptionSubCommand{
				Name:        "info",
				Description: "Returns information about a mnemonic.",
				Options:     opts,
			},
		},
	},
	Handler: func(event *events.ApplicationCommandInteractionCreate, user db.UserEntry) error {
		data := event.SlashCommandInteractionData()

		refreshCredentials, err := blast.RefreshTokenLogin(consts.FORTNITE_PC_CLIENT_ID, consts.FORTNITE_PC_CLIENT_SECRET, user.Accounts[user.SelectedAccount].RefreshToken)
		if err != nil {
			return err
		}

		log.Println(refreshCredentials)

		subCommandName := *data.SubCommandName
		topLevelSubCommandName := subCommandName

		groupNamePointer := data.SubCommandGroupName
		if groupNamePointer != nil {
			topLevelSubCommandName = *data.SubCommandGroupName
		}

		switch topLevelSubCommandName {
		case "favorites":
			switch subCommandName {
			case "list":
				return fmt.Errorf("unknown subcommand")
			case "add":
				err := blast.AddFavoriteMnemonic(refreshCredentials, data.String("mnemonic"))
				if err != nil {
					return err
				}

				return fmt.Errorf("added favorite")
			case "remove":
				err := blast.RemoveFavoriteMnemonic(refreshCredentials, data.String("mnemonic"))
				if err != nil {
					return err
				}

				return fmt.Errorf("removed favorite")
			default:
				return fmt.Errorf("unknown subcommand")
			}
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
