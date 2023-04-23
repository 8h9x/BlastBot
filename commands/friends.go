package commands

import (
	"blast/api"
	"blast/db"
	"encoding/hex"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var friends = discord.SlashCommandCreate{
	Name:        "friends",
	Description: "template command",
	Options: []discord.ApplicationCommandOption{
		// discord.ApplicationCommandOptionSubCommand{
		// 	Name:        "info",
		// 	Description: "Retrieves general account information.",
		// },
		// discord.ApplicationCommandOptionSubCommand{
		// 	Name:        "switch",
		// 	Description: "Switches your selected epic games account.",
		// },
		discord.ApplicationCommandOptionSubCommand{
			Name:        "add",
			Description: "Add a friend.",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:         "user",
					Description:  "friends.add.user.autocomplete_string",
					Autocomplete: true,
					Required:     true,
				},
			},
		},
	},
}

var FriendsAdd = Command{
	Handler: func(event *handler.CommandEvent, blast api.EpicClient, user db.UserEntry, credentials api.UserCredentialsResponse, data discord.SlashCommandInteractionData) error {
		userID := ParseUserID(blast, credentials, data.String("user"))

		err := blast.AddFriend(credentials, userID)
		if err != nil {
			return err
		}

		_, err = event.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
			SetContentf(userID).
			AddActionRow(discord.NewPrimaryButton("test", "test_button")).
			Build(),
		)

		return err
	},
	LoginRequired:     true,
	EphemeralResponse: false,
}

func ParseUserID(blast api.EpicClient, credentials api.UserCredentialsResponse, username string) string {
	if len(username) == 32 && IsHex(username) {
		return username
	}

	platform := "epic"
	splitUser := strings.Split(username, ":")

	if len(splitUser) == 2 {
		platform = splitUser[0]
		username = splitUser[1]
	}

	if !Contains(validPlatforms, platform) {
		return ""
	}

	if platform == "epic" {
		user, err := blast.FetchUserByDisplayName(credentials, username)
		if err != nil {
			return ""
		}

		return user.ID
	}

	user, err := blast.FetchUserByExternalDisplayName(credentials, username, platform)
	if err != nil {
		return ""
	}

	return user[0].ID
}

func IsHex(s string) bool {
	_, err := hex.DecodeString(s)
	return err == nil
}

var validPlatforms = []string{"epic", "psn", "xbl", "steam"}

func Contains(arr []string, str string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}

// https://user-search-service-prod.ol.epicgames.com/api/v1/search/{AccountID}?prefix={displayName}&platform{epic | psn | xbl | steam }
