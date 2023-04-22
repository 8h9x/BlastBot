package commands

import (
	"blast/api"
	"blast/api/consts"
	"blast/db"
	"context"
	"fmt"
	"log"
	"time"

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
		if data.Bool("bulk") {
			col := db.GetCollection("users")

			_, err := col.UpdateOne(context.Background(), bson.M{"discordId": event.User().ID}, bson.M{"$set": bson.M{"accounts": bson.M{}}})
			if err != nil {
				return err
			}

			_, err = col.UpdateOne(context.Background(), bson.M{"discordId": event.User().ID}, bson.M{"$set": bson.M{"selectedAccount": 0}})
			if err != nil {
				return err
			}

			embed := discord.NewEmbedBuilder().
				SetColor(0xFB5A32).
				SetTimestamp(time.Now()).
				SetAuthorIcon(*event.User().AvatarURL(discord.WithFormat(discord.ImageFormatPNG))).
				SetAuthorNamef("Saved accounts bulk removed from %s", event.User().Username).
				SetDescription("You have successfully logged out of all of your saved accounts.\nYou now have (0/25) saved accounts.").
				Build()

			_, err = event.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
				SetEmbeds(embed).
				ClearContent().
				Build(),
			)
			return err
		}

		logOutOptions := make([]discord.StringSelectMenuOption, 0)

		for _, account := range user.Accounts {
			refreshCredentials, err := blast.RefreshTokenLogin(consts.FORTNITE_PC_CLIENT_ID, consts.FORTNITE_PC_CLIENT_SECRET, account.RefreshToken)
			if err != nil {
				if err.(*api.RequestError).Raw.ErrorCode == "errors.com.epicgames.account.auth_token.invalid_refresh_token" {
					log.Println(account.AccountID)
					col := db.GetCollection("users")
					_, err = col.UpdateOne(context.Background(), bson.M{"discordId": event.User().ID}, bson.M{"$pull": bson.M{"accounts": bson.M{"accountId": account.AccountID}}})
					if err != nil {
						return err
					}

					continue
				}
				return err
			}

			accountInfo, err := blast.FetchAccountInformation(refreshCredentials)
			if err != nil {
				return err
			}

			logOutOptions = append(logOutOptions, discord.NewStringSelectMenuOption(accountInfo.DisplayName, fmt.Sprintf("%s:%s", account.AccountID, accountInfo.DisplayName)))
		}

		_, err := event.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
			AddActionRow(discord.NewStringSelectMenu("logout_account_select", "Select an account to log out of", logOutOptions...)).
			AddActionRow(discord.NewDangerButton("Cancel", "cancel").WithEmoji(discord.ComponentEmoji{
				Name: "x_",
				ID:   1096630553689739385,
			})).
			Build(),
		)
		return err

		// return fmt.Errorf("successfully logged out of %s", account.DisplayName)
	},
	LoginRequired:     true,
	EphemeralResponse: false,
}
