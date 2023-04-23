package components

import (
	"blast/api"
	"blast/api/consts"
	"blast/db"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"go.mongodb.org/mongo-driver/bson"
)

// vaid platforms: epic | psn | xbl | steam

func UserSearchAutocomplete(e *handler.AutocompleteEvent) error {
	username := e.Data.String("user")
	platform := "epic"

	splitUser := strings.Split(username, ":")

	if len(splitUser) == 2 {
		platform = splitUser[0]
		username = splitUser[1]
	}

	if !contains(validPlatforms, platform) {
		return nil
	}

	if username == "" || len(username) < 2 || len(username) > 32 {
		return nil
	}

	if len(username) == 32 && IsHex(username) {
		return nil
	}

	user, err := db.Fetch[db.UserEntry]("users", bson.M{"discordId": e.User().ID})
	if err != nil {
		return err
	}

	blast := api.New()

	credentials, err := blast.RefreshTokenLogin(consts.FORTNITE_PC_CLIENT_ID, consts.FORTNITE_PC_CLIENT_SECRET, user.Accounts[user.SelectedAccount].RefreshToken)
	if err != nil {
		return err
	}

	users, err := blast.UserSearch(credentials, username, platform)
	if err != nil {
		return err
	}

	if len(users) == 0 {
		return nil
	}

	var res []discord.AutocompleteChoice

	for _, user := range users {
		if len(res) == 24 {
			break
		}

		displayName := user.Matches[0].Value
		if user.Matches[0].Platform != "epic" {
			displayName = fmt.Sprintf("%s:%s", user.Matches[0].Platform, displayName)
		}

		res = append(res, discord.AutocompleteChoiceString{
			Name:  displayName,
			Value: user.AccountID,
		})
	}

	err = e.Result(res)
	if err != nil {
		return err
	}

	return err
}

func IsHex(s string) bool {
	_, err := hex.DecodeString(s)
	return err == nil
}

var validPlatforms = []string{"epic", "psn", "xbl", "steam"}

func contains(arr []string, str string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}
