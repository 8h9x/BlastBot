package components

import (
	"blast/db"
	"context"
	"strings"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"go.mongodb.org/mongo-driver/bson"
)

func LogoutAccountSelect(e *handler.ComponentEvent) error {
	user := e.User()

	data := strings.Split(e.StringSelectMenuInteractionData().Values[0], ":")
	accountID, displayName := data[0], data[1]

	col := db.GetCollection("users")

	_, err := col.UpdateOne(context.Background(), bson.M{"discordId": user.ID.String()}, bson.M{"$pull": bson.M{"accounts": bson.M{"accountId": accountID}}})
	if err != nil {
		return err
	}

	userEntry, err := db.Fetch[db.UserEntry]("users", bson.M{"discordId": user.ID.String()})
	if err != nil {
		return err
	}

	return e.UpdateMessage(discord.NewMessageUpdateBuilder().
		ClearContent().
		ClearContainerComponents().
		SetEmbeds(discord.NewEmbedBuilder().
			SetColor(0xFB5A32).
			SetTimestamp(time.Now()).
			SetAuthorIcon(*user.AvatarURL(discord.WithFormat(discord.ImageFormatPNG))).
			SetAuthorNamef("Saved account removed from %s", user.Username).
			SetDescriptionf("Successfully logged out of **%s**\nYou now have (%d/25) saved accounts.", displayName, len(userEntry.Accounts)).
			Build()).
		Build(),
	)
}
