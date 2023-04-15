package commands

import (
	"blast/api"
	"blast/api/consts"
	"blast/db"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var mcp = discord.SlashCommandCreate{
	Name:        "mcp",
	Description: "Performs a custom MCP request.",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:        "operation",
			Description: "The name of the operation to perform (ex. QueryProfile)", // unsure yet if case-sensitive.
			Required:    true,
		},
		discord.ApplicationCommandOptionString{
			Name:        "profile",
			Description: "The profile to perform the operation on.",
			Choices: []discord.ApplicationCommandOptionChoiceString{
				{Name: "athena (Battle Royale)", Value: "athena"},
				{Name: "creative", Value: "creative"},
				{Name: "campaign (Save The World)", Value: "campaign"},
				{Name: "common_public (Public Data)", Value: "common_public"},
				{Name: "collections (BR Collections)", Value: "collections"},
				{Name: "common_core (Main Profile)", Value: "common_core"},
				{Name: "metadata (STW Homebase)", Value: "metadata"},
				{Name: "collection_book_people0 (STW Collection Book People)", Value: "collection_book_people0"},
				{Name: "collection_book_schematics0 (STW Collection Book Schematics)", Value: "collection_book_schematics0"},
				{Name: "outpost0 (STW Storage)", Value: "outpost0"},
				{Name: "theater0 (STW Backpack)", Value: "theater0"},
				{Name: "theater1 (STW Event Backpack)", Value: "theater1"},
				{Name: "theater2 (STW Ventures Data)", Value: "theater2"},
				{Name: "recycle_bin (STW Recylce Bin)", Value: "recycle_bin"},
			},
			Required: true,
		},
		discord.ApplicationCommandOptionAttachment{
			Name:        "body",
			Description: "The JSON body to be sent with the request (if ommitted, '{}' will be used).",
		},
	},
}

var MCP = Command{
	Handler: func(event *handler.CommandEvent, blast api.EpicClient, user db.UserEntry, credentials api.UserCredentialsResponse, data discord.SlashCommandInteractionData) error {
		refreshCredentials, err := blast.RefreshTokenLogin(consts.FORTNITE_PC_CLIENT_ID, consts.FORTNITE_PC_CLIENT_SECRET, user.Accounts[user.SelectedAccount].RefreshToken)
		if err != nil {
			return err
		}

		operation, profileId, customBodyURL, customBody := data.String("operation"), data.String("profile"), data.Attachment("body").URL, "{}"

		if customBodyURL != "" {
			body, err := fetchJSON(customBodyURL)
			if err != nil {
				return err
			}

			customBody = body
		}

		res, err := blast.ProfileOperationStr(refreshCredentials, operation, profileId, customBody)
		if err != nil {
			return err
		}

		_, err = event.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
			AddFile(fmt.Sprintf("%s.%s.blast.json", strings.ToLower(operation), profileId), "Profile operation response.", res).
			Build(),
		)
		if err != nil {
			return err
		}

		res.Close()

		return nil
	},
	LoginRequired:     true,
	EphemeralResponse: false,
}

func fetchJSON(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
