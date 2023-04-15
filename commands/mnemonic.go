package commands

import (
	"blast/api"
	"blast/db"
	"fmt"
	"log"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var opts = []discord.ApplicationCommandOption{
	discord.ApplicationCommandOptionString{
		Name:        "mnemonic",
		Description: "Playlist id or creative map code.",
		Required:    true,
	},
}

var mnemonic = discord.SlashCommandCreate{
	Name:        "mnemonic",
	Description: "Mnemonic related commands.",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionSubCommandGroup{
			Name:        "favorites",
			Description: "Mnemonic favorites related commands.",
			Options: []discord.ApplicationCommandOptionSubCommand{
				{
					Name:        "add",
					Description: "Adds a mnemonic to your favorites list.",
					Options:     opts,
				},
				{
					Name:        "list",
					Description: "Returns your discovery favorites list.",
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
}

var MnemonicFavoritesAdd = Command{
	Handler: func(event *handler.CommandEvent, blast api.EpicClient, user db.UserEntry, credentials api.UserCredentialsResponse, data discord.SlashCommandInteractionData) error {
		err := blast.AddFavoriteMnemonic(credentials, data.String("mnemonic"))
		if err != nil {
			return err
		}

		return fmt.Errorf("added favorite")
	},
	LoginRequired:     true,
	EphemeralResponse: false,
}

var MnemonicFavoritesList = Command{
	Handler: func(event *handler.CommandEvent, blast api.EpicClient, user db.UserEntry, credentials api.UserCredentialsResponse, data discord.SlashCommandInteractionData) error {
		return fmt.Errorf("soon")
	},
	LoginRequired:     true,
	EphemeralResponse: false,
}

var MnemonicFavoritesRemove = Command{
	Handler: func(event *handler.CommandEvent, blast api.EpicClient, user db.UserEntry, credentials api.UserCredentialsResponse, data discord.SlashCommandInteractionData) error {
		err := blast.RemoveFavoriteMnemonic(credentials, data.String("mnemonic"))
		if err != nil {
			return err
		}

		return fmt.Errorf("removed favorite")
	},
	LoginRequired:     true,
	EphemeralResponse: false,
}

var MnemonicInfo = Command{
	Handler: func(event *handler.CommandEvent, blast api.EpicClient, user db.UserEntry, credentials api.UserCredentialsResponse, data discord.SlashCommandInteractionData) error {
		res, err := blast.FetchMnemonicInfo(credentials, data.String("mnemonic"))
		if err != nil {
			return err
		}

		log.Println(res)

		return fmt.Errorf("data logged to console")
	},
	LoginRequired:     true,
	EphemeralResponse: false,
}
