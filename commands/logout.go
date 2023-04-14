package commands

import (
	"blast/api/consts"
	"blast/db"
	"context"
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"go.mongodb.org/mongo-driver/bson"
)

var logout = Command{
	Create: discord.SlashCommandCreate{
		Name:        "logout",
		Description: "Log out of one (or all) of your saved Epic Games accounts.",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionBool{
				Name:        "bulk",
				Description: "Log out of all of your saved Epic Games accounts.",
			},
		},
	},
	Handler: func(event *events.ApplicationCommandInteractionCreate, user db.UserEntry) error {
		if len(user.Accounts) == 0 {
			return event.CreateMessage(discord.NewMessageCreateBuilder().
				SetContent("You don't have any accounts saved!").
				Build(),
			)
		}

		// embed := discord.NewEmbedBuilder().
		// 	SetColor(0xFB5A32).
		// 	SetTimestamp(time.Now()).
		// 	SetTitle("Log out of an account <a:blastrocket:1094632950395588608>").
		// 	SetDescription("Select an account to log out of.").
		// 	Build()

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

		refreshCredentials, err := blast.RefreshTokenLogin(consts.FORTNITE_PC_CLIENT_ID, consts.FORTNITE_PC_CLIENT_SECRET, user.Accounts[user.SelectedAccount].RefreshToken)
		if err != nil {
			return err
		}

		account, err := blast.FetchAccountInformation(refreshCredentials)
		if err != nil {
			return err
		}

		col := db.GetCollection("users")

		_, err = col.UpdateOne(context.Background(), bson.M{"discordId": event.User().ID.String()}, bson.M{"$pull": bson.M{"accounts": bson.M{"accountId": user.Accounts[user.SelectedAccount].AccountID}}})
		if err != nil {
			return err
		}

		// return fmt.Errorf("successfully logged out of %s", account.DisplayName)

		// return nil

		_, err = event.Client().Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.NewMessageUpdateBuilder().
			// SetEmbeds(embed).
			AddActionRow(
				discord.NewStringSelectMenu("account", "Select an account to log out of").
					AddOptions(logOutOptions...),
				// discord.NewDangerButton("Cancel", "cancel-logout"),
			).
			Build(),
		)
		if err != nil {
			return err
		}

		return fmt.Errorf("successfully logged out of %s", account.DisplayName)

		// embed := discord.NewEmbedBuilder().
		// 	// SetAuthorIcon(event.). // TODO set author icon to bot user avatar
		// 	SetColor(0xFB5A32).
		// 	SetTimestamp(time.Now()).
		// 	SetTitle("Add a new account <a:blastrocket:1094632950395588608>").
		// 	SetDescriptionf("**Login Instructions:**\n**1.** Click the `Login` button below.\n**2.** Click the `Confirm` button on the epic games page.\n**3.** Wait a few seconds for the bot to process login.\n\n***This interaction will timeout <t:%d:R>.***", expires).
		// 	Build()

		// err = event.CreateMessage(discord.NewMessageCreateBuilder().
		// 	SetEmbeds(embed).
		// 	AddActionRow(
		// 		discord.NewLinkButton("Login", deviceAuthorization.VerificationUriComplete),
		// 		discord.NewDangerButton("Cancel", "cancel-login"),
		// 	).
		// 	Build(),
		// )
		// if err != nil {
		// 	return err
		// }

		// return fmt.Errorf("not implemented")
	},
	LoginRequired: true,
}
