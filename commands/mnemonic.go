package commands

import (
	"blast/api"
	"blast/db"
	"fmt"
	"time"

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
		mnemonic := data.String("mnemonic")

		err := blast.AddFavoriteMnemonic(credentials, mnemonic)
		if err != nil {
			return err
		}

		info, err := blast.FetchMnemonicInfo(credentials, mnemonic)
		if err != nil {
			return err
		}

		_, err = event.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
			SetEmbeds(discord.NewEmbedBuilder().
				SetTitlef("Mnemonic has been favorited!").
				SetDescriptionf("`%s | %s`", info.Metadata.Title, mnemonic).
				SetColor(0xFB5A32).
				SetTimestamp(time.Now()).
				Build(),
			).
			Build(),
		)
		return err
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
		mnemonic := data.String("mnemonic")

		err := blast.RemoveFavoriteMnemonic(credentials, mnemonic)
		if err != nil {
			return err
		}

		info, err := blast.FetchMnemonicInfo(credentials, mnemonic)
		if err != nil {
			return err
		}

		_, err = event.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
			SetEmbeds(discord.NewEmbedBuilder().
				SetTitlef("Mnemonic has been unfavorited!").
				SetDescriptionf("`%s | %s`", info.Metadata.Title, mnemonic).
				SetColor(0xFB5A32).
				SetTimestamp(time.Now()).
				Build(),
			).
			Build(),
		)
		return err
	},
	LoginRequired:     true,
	EphemeralResponse: false,
}

var MnemonicInfo = Command{
	Handler: func(event *handler.CommandEvent, blast api.EpicClient, user db.UserEntry, credentials api.UserCredentialsResponse, data discord.SlashCommandInteractionData) error {
		mnemonic := data.String("mnemonic")

		info, err := blast.FetchMnemonicInfo(credentials, mnemonic)
		if err != nil {
			return err
		}

		imageURL := info.Metadata.ImageURL
		if imageURL == "" {
			imageURL = info.Metadata.GeneratedIslandUrlsOld.URL
		}

		_, err = event.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
			SetEmbeds(discord.NewEmbedBuilder().
				SetTitlef("%s | %s", info.Metadata.Title, mnemonic).
				SetDescription(info.Metadata.Tagline).
				SetThumbnail(imageURL).
				SetColor(0xFB5A32).
				SetTimestamp(time.Now()).
				AddField("Author", info.CreatorName, true).
				AddField("Version", fmt.Sprint(info.Version), true).
				AddField("Link Type", info.LinkType, true).
				AddField("Publish Date", info.Published.String(), true).
				AddField("XP Calibration Phase", info.Metadata.DynamicXp.CalibrationPhase, true).
				AddField("Active", BoolToEmoji(info.Active), true).
				Build(),
			).
			Build(),
		)
		return err
	},
	LoginRequired:     true,
	EphemeralResponse: false,
}
