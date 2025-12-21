package mcp

import (
	"fmt"
	"io"
	"net/http"
	"time"

    "github.com/8h9x/BlastBot/internal/sessions"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var Definition = discord.SlashCommandCreate{
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

func Handler(event *handler.CommandEvent) error {
	discordId := event.User().ID
	data := event.SlashCommandInteractionData()
	operation, profileId, customBodyURL, customBody := data.String("operation"), data.String("profile"), data.Attachment("body").URL, "{}"

	if customBodyURL != "" {
		body, err := fetchJSON(customBodyURL)
		if err != nil {
			return err
		}

		customBody = body
	}

	session, err := sessions.GetSessionForUser(discordId)
	if err != nil {
		return fmt.Errorf("unable to get session for user (%s): %s", discordId, err)
	}

	res, err := session.ComposeProfileOperation(operation, profileId, customBody)
	if err != nil {
		return err
	}

	// profile, err := request.ResponseParser[vinderman.Profile[vinderman.AthenaProfileStats, []interface{}]](res)
	// if err != nil {
	// 	return err
	// }

	err = event.CreateMessage(discord.NewMessageCreateBuilder().
		AddFile(fmt.Sprintf("profile_%s_%d.json", profileId, time.Now().Unix()), "Profile operation response.", res.Body).
		Build(),
	)
	if err != nil {
		return err
	}

	return nil
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
