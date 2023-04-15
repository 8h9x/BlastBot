package commands

import (
	"blast/api"
	"blast/api/consts"
	"blast/db"
	"context"
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"go.mongodb.org/mongo-driver/bson"
)

var logout = discord.SlashCommandCreate{
	Name:        "logout",
	Description: "Log out of one (or all) of your saved Epic Games accounts.",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionBool{
			Name:        "bulk",
			Description: "Log out of all of your saved Epic Games accounts.",
		},
	},
}

var Logout = Command{
	Handler: func(event *handler.CommandEvent, blast api.EpicClient, user db.UserEntry, credentials api.UserCredentialsResponse, data discord.SlashCommandInteractionData) error {
		logOutOptions := make([]discord.StringSelectMenuOption, 0)

		for i, account := range user.Accounts {
			refreshCredentials, err := blast.RefreshTokenLogin(consts.FORTNITE_PC_CLIENT_ID, consts.FORTNITE_PC_CLIENT_SECRET, account.RefreshToken)
			if err != nil {
				return err
			}

			accountInfo, err := blast.FetchAccountInformation(refreshCredentials)
			if err != nil {
				return err
			}

			logOutOptions = append(logOutOptions, discord.NewStringSelectMenuOption(accountInfo.DisplayName, account.AccountID+fmt.Sprint(i)))
		}

		account, err := blast.FetchAccountInformation(credentials)
		if err != nil {
			return err
		}

		col := db.GetCollection("users")

		_, err = col.UpdateOne(context.Background(), bson.M{"discordId": event.User().ID.String()}, bson.M{"$pull": bson.M{"accounts": bson.M{"accountId": user.Accounts[user.SelectedAccount].AccountID}}})
		if err != nil {
			return err
		}

		_, err = event.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
			// SetEmbeds(embed).
			AddActionRow(discord.NewStringSelectMenu("account", "Select an account to log out of", logOutOptions...)).
			AddActionRow(discord.NewDangerButton("Cancel", "cancel-logout")).
			Build(),
		)
		if err != nil {
			return err
		}

		return fmt.Errorf("successfully logged out of %s", account.DisplayName)
	},
	LoginRequired:     true,
	EphemeralResponse: false,
}
