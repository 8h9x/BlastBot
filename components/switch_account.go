package components

import (
	"blast/api"
	"blast/api/consts"
	"blast/db"
	"context"
	"strconv"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"go.mongodb.org/mongo-driver/bson"
)

func SwitchAccountSelect(e *handler.ComponentEvent) error {
	user := e.User()

	accountIndex, err := strconv.Atoi(e.StringSelectMenuInteractionData().Values[0])
	if err != nil {
		return err
	}

	col := db.GetCollection("users")

	_, err = col.UpdateOne(context.Background(), bson.M{"discordId": user.ID}, bson.M{"$set": bson.M{"selectedAccount": accountIndex}})
	if err != nil {
		return err
	}

	userEntry, err := db.Fetch[db.UserEntry]("users", bson.M{"discordId": user.ID})
	if err != nil {
		return err
	}

	blast := api.New()

	credentials, err := blast.RefreshTokenLogin(consts.FORTNITE_PC_CLIENT_ID, consts.FORTNITE_PC_CLIENT_SECRET, userEntry.Accounts[userEntry.SelectedAccount].RefreshToken)
	if err != nil {
		return err
	}

	account, err := blast.FetchMyAccountInfo(credentials)
	if err != nil {
		return err
	}

	avatarURL, err := blast.FetchAvatarURL(credentials)
	if err != nil {
		return err
	}

	return e.UpdateMessage(discord.NewMessageUpdateBuilder().
		ClearContent().
		ClearContainerComponents().
		SetEmbeds(discord.NewEmbedBuilder().
			SetColor(0xFB5A32).
			SetTimestamp(time.Now()).
			SetAuthorIcon(avatarURL).
			SetAuthorNamef("Selected account switched to %s", account.DisplayName).
			Build()).
		Build(),
	)
}
