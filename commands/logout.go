package commands

import (
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
	Handler: func(event *events.ApplicationCommandInteractionCreate) error {
		user, err := db.Fetch[db.UserEntry]("users", bson.M{"discordId": event.User().ID.String()})
		if err != nil {
			return err
		}

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

		for _, account := range user.Accounts {
			logOutOptions = append(logOutOptions, discord.NewStringSelectMenuOption(account.AccountID, account.AccountID))
		}

		col := db.GetCollection("users")

		_, err = col.UpdateOne(context.Background(), bson.M{"discordId": event.User().ID.String()}, bson.M{"$pull": bson.M{"accounts": bson.M{"accountId": user.Accounts[user.SelectedAccount].AccountID}}})
		if err != nil {
			return err
		}

		return fmt.Errorf("successfully logged out of %s", user.Accounts[user.SelectedAccount].AccountID)

		// return nil

		// err = event.CreateMessage(discord.NewMessageCreateBuilder().
		// 	// SetEmbeds(embed).
		// 	AddActionRow(
		// 		discord.NewStringSelectMenu("account", "Select an account to log out of").
		// 			AddOptions(logOutOptions...),
		// 		// discord.NewDangerButton("Cancel", "cancel-logout"),
		// 	).
		// 	Build(),
		// )
		// if err != nil {
		// 	return err
		// }

		// return nil

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
