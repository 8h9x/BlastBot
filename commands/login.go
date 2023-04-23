package commands

import (
	"blast/api"
	"blast/api/consts"
	"blast/db"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
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
		if len(user.Accounts) >= 25 {
			embed := discord.NewEmbedBuilder().
				SetColor(0xFB5A32).
				SetTimestamp(time.Now()).
				SetAuthorIcon(*event.User().AvatarURL(discord.WithFormat(discord.ImageFormatPNG))).
				SetAuthorNamef("Saved account limit reached for %s", event.User().Username).
				SetDescription("You have reached the maximum amount of saved accounts.\nRemove one using the `/logout` command.").
				Build()

			_, err := event.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
				SetEmbeds(embed).
				ClearContent().
				Build(),
			)
			return err
		}

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
				discord.NewLinkButton("Login", deviceAuthorization.VerificationUriComplete).WithEmoji(discord.ComponentEmoji{
					Name: "add",
					ID:   1096633864136446143,
				}),
				discord.NewDangerButton("Cancel", "cancel").WithEmoji(discord.ComponentEmoji{
					Name: "x_",
					ID:   1096630553689739385,
				}),
			).
			Build(),
		)
		if err != nil {
			return err
		}

		deviceCodeCredentials, err := WaitForDeviceCodeAcceptEOS(blast, consts.UEFN_CLIENT_ID, consts.UEFN_CLIENT_SECRET, deviceAuthorization.DeviceCode, 0)
		if err != nil {
			if err.Error() == "errors.blast.timed_out" {
				return event.DeleteInteractionResponse()
			}

			return err
		}

		for _, account := range user.Accounts {
			if account.AccountID == deviceCodeCredentials.AccountID {
				embed := discord.NewEmbedBuilder().
					SetColor(0xFB5A32).
					SetTimestamp(time.Now()).
					SetAuthorIcon(*event.User().AvatarURL(discord.WithFormat(discord.ImageFormatPNG))).
					SetAuthorNamef("Account already linked to %s", event.User().Username).
					SetDescription("You have already saved this account.\nTo remove it, use the `/logout` command.").
					Build()

				_, err := event.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
					SetEmbeds(embed).
					ClearContent().
					Build(),
				)
				return err
			}
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

		userId := event.User().ID

		_, err = db.Fetch[db.UserEntry]("users", bson.M{"discordId": userId})

		col := db.GetCollection("users")

		userEntry := db.UserEntry{
			ID:        primitive.NewObjectID().Hex(),
			DiscordID: userId,
			Accounts: []db.EpicAccountEntry{
				{
					AccountID:        exchangeCredentials.AccountID,
					RefreshToken:     refreshPayload.Jti,
					RefreshExpiresAt: exchangeCredentials.RefreshExpiresAt,
					ClientId:         exchangeCredentials.ClientId,
					Flags:            db.USER,
				},
			},
			SelectedAccount: 0,
			BulkFlags:       db.USER,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		if err == nil { // user exists
			_, err := col.UpdateOne(context.Background(), bson.M{"discordId": userId}, bson.M{"$push": bson.M{"accounts": bson.M{
				"accountId":        exchangeCredentials.AccountID,
				"refreshToken":     refreshPayload.Jti,
				"refreshExpiresAt": exchangeCredentials.RefreshExpiresAt,
				"clientId":         exchangeCredentials.ClientId,
				"flags":            db.AUTODAILY,
			}}}, options.Update().SetUpsert(true))
			if err != nil {
				return err
			}

			_, err = col.UpdateOne(context.Background(), bson.M{"discordId": userId}, bson.M{"$set": bson.M{"selectedAccount": len(user.Accounts)}})
			if err != nil {
				return err
			}
		} else {
			_, err = col.InsertOne(context.Background(), userEntry)
			if err != nil {
				return err
			}
		}

		refreshCredentials, err := blast.RefreshTokenLogin(consts.FORTNITE_PC_CLIENT_ID, consts.FORTNITE_PC_CLIENT_SECRET, refreshPayload.Jti)
		if err != nil {
			return err
		}

		accountInfo, err := blast.FetchMyAccountInfo(refreshCredentials)
		if err != nil {
			return err
		}

		avatarURL, err := blast.FetchAvatarURL(refreshCredentials)
		if err != nil {
			return err
		}

		_, err = event.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
			ClearContent().
			ClearContainerComponents().
			SetEmbeds(discord.NewEmbedBuilder().
				SetColor(0xFB5A32).
				SetTimestamp(time.Now()).
				SetAuthorIcon(avatarURL).
				SetAuthorNamef("New account saved for %s", event.User().Username).
				SetDescriptionf("Successfully logged into **%s**\nYou now have (%d/25) saved accounts.", accountInfo.DisplayName, len(user.Accounts)+1).
				Build()).
			Build(),
		)

		// _, err = event.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().SetContentf("Logged in as `%s`!", deviceCodeCredentials.AccountID).Build())
		// if err != nil {
		// 	return err
		// }

		return err
	},
	LoginRequired:     false,
	EphemeralResponse: false,
}

func base64Decode(s string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(s)
}

// redo this here because goroutines are so annoying to work with when trying to cancel early
func WaitForDeviceCodeAcceptEOS(c api.EpicClient, clientId string, clientSecret string, deviceCode string, runs int) (api.UserCredentialsResponseEOS, error) {
	if runs == 12 { // after 120 seconds
		return api.UserCredentialsResponseEOS{}, errors.New("errors.blast.timed_out")
	}

	deviceCodeCredentials, err := c.DeviceCodeLoginEOS(clientId, clientSecret, deviceCode)

	if err != nil {
		if err.(*api.RequestError).Raw.ErrorCode == "errors.com.epicgames.account.oauth.authorization_pending" {
			time.Sleep(10 * time.Second)
			return WaitForDeviceCodeAcceptEOS(c, clientId, clientSecret, deviceCode, runs+1)
		}

		return api.UserCredentialsResponseEOS{}, err
	}

	return deviceCodeCredentials, nil
}
