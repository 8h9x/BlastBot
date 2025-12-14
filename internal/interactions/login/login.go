package login

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/8h9x/BlastBot/database/internal/database"
	"github.com/8h9x/BlastBot/database/internal/manager/sessions"
	"github.com/8h9x/fortgo"
	"github.com/8h9x/fortgo/auth"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	CHECK_INTERVAL = time.Second * 10
	CHECK_TIMEOUT  = time.Minute * 2
)

var Definition = discord.SlashCommandCreate{
	Name:        "login",
	Description: "Login and grant Blast! access to your epic account. (alias for /accounts add)",
}

func Handler(event *handler.CommandEvent) error {
	httpClient := &http.Client{}

	clientCredentials, err := auth.Authenticate(httpClient, auth.FortnitePS4USClient, auth.PayloadClientCredentials{}, false)
	if err != nil {
		return err
	}

	deviceAuthorization, err := auth.GetDeviceCode(httpClient, clientCredentials)
	if err != nil {
		return err
	}

	err = event.CreateMessage(discord.MessageCreate{
		Content: fmt.Sprintf("Visit this URL: %s\nThen press 'Confirm' on the epic games login page.\n", deviceAuthorization.VerificationURIComplete),
	})
	if err != nil {
		return err
	}

	credentials, err := waitForDeviceCodeConfirm(httpClient, deviceAuthorization.DeviceCode, CHECK_INTERVAL, CHECK_TIMEOUT)
	if err != nil {
		return err
	}

	_, err = fortgo.NewClient(httpClient, credentials)
	if err != nil {
		return err
	}

	userId := event.User().ID

	userEntry := database.User{
		ID:        bson.NewObjectID().Hex(),
		DiscordID: userId,
		Accounts: []database.EpicAccount{
			{
				AccountID:        credentials.AccountID,
				RefreshToken:     credentials.RefreshToken,
				RefreshExpiresAt: credentials.RefreshExpiresAt,
				CreatedClientID:  credentials.ClientID,
				Flags:            database.USER,
			},
		},
		SelectedEpicAccountId: credentials.AccountID,
		BulkFlags:             database.USER,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}

	col := database.GetCollection("users")

	_, err = database.Fetch[database.User]("users", bson.M{"discordId": userId})
	if err == nil { // user exists
		_, err := col.UpdateOne(context.Background(), bson.M{"discordId": userId}, bson.M{"$push": bson.M{"accounts": bson.M{
			"accountId":        credentials.AccountID,
			"refreshToken":     credentials.RefreshToken,
			"refreshExpiresAt": credentials.RefreshExpiresAt,
			"clientId":         credentials.ClientID,
			"flags":            database.USER,
		}}}, options.UpdateOne().SetUpsert(true))
		if err != nil {
			return err
		}

		_, err = col.UpdateOne(context.Background(), bson.M{"discordId": userId}, bson.M{"$set": bson.M{"selectedEpicAccountId": credentials.AccountID}})
		if err != nil {
			return err
		}
	} else {
		_, err = col.InsertOne(context.Background(), userEntry)
		if err != nil {
			return err
		}
	}

	session, err := sessions.CreateSession(httpClient, credentials)
	if err != nil {
		return err
	}

	_, err = event.CreateFollowupMessage(discord.NewMessageCreateBuilder().
		ClearContainerComponents().
		SetEmbeds(discord.NewEmbedBuilder().
			SetColor(0xFB5A32).
			SetTimestamp(time.Now()).
			SetAuthorNamef("New account saved for %s", event.User().Username).
			SetDescriptionf("Successfully logged in using client: **%s**\nYou now have (%d/25) saved accounts.", session.ClientID, len(userEntry.Accounts)+1).
			Build()).
		Build(),
	)

	return err
}

func waitForDeviceCodeConfirm(httpClient *http.Client, deviceCode string, interval, timeout time.Duration) (auth.TokenResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	payload := auth.PayloadDeviceCode{
		DeviceCode: deviceCode,
	}

	for {
		select {
		case <-ctx.Done():
			return auth.TokenResponse{}, ctx.Err()
		case <-ticker.C:
			credentials, err := auth.Authenticate(httpClient, auth.FortnitePS4USClient, payload, true)
			if err != nil {
				continue
			}
			return credentials, nil
		}
	}
}
