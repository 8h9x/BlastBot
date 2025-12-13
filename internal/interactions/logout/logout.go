package logout

import (
	"context"

	"github.com/8h9x/BlastBot/database/internal/database"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var Definition = discord.SlashCommandCreate{
	Name:        "logout",
	Description: "Removes currently selected account from the bot. (alias for /accounts remove {current})",
}

func Handler(event *handler.CommandEvent) error {
	discordId := event.User().ID
	col := database.GetCollection("users")

	result, err := database.Fetch[database.User]("users", bson.M{"discordId": discordId})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = event.CreateMessage(discord.MessageCreate{
				Content: "You do not have any saved accounts.",
			})
			if err != nil {
				return err
			}
		}
		return err
	}

	if len(result.Accounts) == 0 {
		err = event.CreateMessage(discord.MessageCreate{
			Content: "You do not have any saved accounts.",
		})
		if err != nil {
			return err
		}
	}

	var newSelectedId string = ""
	for _, account := range result.Accounts {
		if account.AccountID != result.SelectedEpicAccountId {
			newSelectedId = account.AccountID
			break
		}
	}

	_, err = col.UpdateOne(
		context.Background(),
		bson.M{"discordId": discordId},
		bson.M{
			"$pull": bson.M{
				"accounts": bson.M{"accountId": result.SelectedEpicAccountId},
			},
			"$set": bson.M{
				"selectedEpicAccountId": newSelectedId,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}
