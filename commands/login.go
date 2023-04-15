package commands

import (
	"blast/api"
	"blast/api/consts"
	"blast/db"
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var login = discord.SlashCommandCreate{
	Name:        "login",
	Description: "Login to your Epic Games account.",
}

var Login = Command{
	Handler: func(event *handler.CommandEvent, blast api.EpicClient, user db.UserEntry, credentials api.UserCredentialsResponse, data discord.SlashCommandInteractionData) error {
		clientCredentials, err := blast.GetClientCredentialsEOS(consts.UEFN_CLIENT_ID, consts.UEFN_CLIENT_SECRET)
		if err != nil {
			return err
		}

		deviceAuthorization, err := blast.GetDeviceAuthorizationEOS(clientCredentials.AccessToken)
		if err != nil {
			log.Fatal(err)
		}

		expires := time.Now().Add(time.Minute * 2).Unix()

		embed := discord.NewEmbedBuilder().
			// SetAuthorIcon(event.). // TODO set author icon to bot user avatar
			SetColor(0xFB5A32).
			SetTimestamp(time.Now()).
			SetTitle("Add a new account <a:rocket:1094632950395588608>").
			SetDescriptionf("**Login Instructions:**\n**1.** Click the `Login` button below.\n**2.** Click the `Confirm` button on the epic games page.\n**3.** Wait a few seconds for the bot to process login.\n\n***This interaction will timeout <t:%d:R>.***", expires).
			Build()

		_, err = event.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
			SetEmbeds(embed).
			AddActionRow(
				discord.NewLinkButton("Login", deviceAuthorization.VerificationUriComplete),
				discord.NewDangerButton("Cancel", "cancel-login"),
			).
			Build(),
		)
		if err != nil {
			return err
		}

		deviceCodeCredentials, err := blast.WaitForDeviceCodeAcceptEOS(consts.UEFN_CLIENT_ID, consts.UEFN_CLIENT_SECRET, deviceAuthorization.DeviceCode)
		if err != nil {
			return err
		}

		exchangeCode, err := blast.GetExchangeCodeEOS(deviceCodeCredentials)
		if err != nil {
			return err
		}

		// decode this refresh jwt and get jti for hex format refresh token that expires after ~6 months (170 days) and can be infinitly refreshed
		exchangeCredentials, err := blast.ExchangeCodeLoginEOS(consts.FORTNITE_PC_CLIENT_ID, consts.FORTNITE_PC_CLIENT_SECRET, exchangeCode.Code)
		if err != nil {
			return err
		}

		decodedRefreshJwtPayload, err := base64Decode(strings.Split(exchangeCredentials.RefreshToken, ".")[1])
		if err != nil {
			return err
		}

		refreshPayload := api.DecodedRefreshTokenJwtPayload{}

		err = json.Unmarshal(decodedRefreshJwtPayload, &refreshPayload)
		if err != nil {
			return err
		}

		userId := event.User().ID.String()

		_, err = db.Fetch[db.UserEntry]("users", bson.M{"discordId": userId})

		col := db.GetCollection("users")

		userEntry := db.UserEntry{
			ID:        primitive.NewObjectID().Hex(),
			DiscordID: userId,
			Accounts: []db.EpicAccountEntry{
				{
					AccountID:        exchangeCredentials.AccountId,
					RefreshToken:     refreshPayload.Jti,
					RefreshExpiresAt: exchangeCredentials.RefreshExpiresAt,
					ClientId:         exchangeCredentials.ClientId,
					Flags: db.AccountFlags{
						AutoDailyClaim: false,
					},
				},
			},
			SelectedAccount: 0,
			BulkFlags: db.AccountFlags{
				AutoDailyClaim: false,
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err == nil { // user exists
			_, err := col.UpdateOne(context.Background(), bson.M{"discordId": userId}, bson.M{"$push": bson.M{"accounts": bson.M{
				"accountId":        exchangeCredentials.AccountId,
				"refreshToken":     refreshPayload.Jti,
				"refreshExpiresAt": exchangeCredentials.RefreshExpiresAt,
				"clientId":         exchangeCredentials.ClientId,
				"flags": bson.M{
					"autoDailyClaim": false,
				},
			}}}, options.Update().SetUpsert(true))
			if err != nil {
				return err
			}
		} else {
			_, err = col.InsertOne(context.Background(), userEntry)
			if err != nil {
				return err
			}
		}

		_, err = event.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().SetContentf("Logged in as `%s`!", deviceCodeCredentials.AccountId).Build())
		if err != nil {
			return err
		}

		return nil
	},
	LoginRequired:     false,
	EphemeralResponse: false,
}

func base64Decode(s string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(s)
}
