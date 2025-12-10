package commands

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/8h9x/BlastBot/database/internal/database"
	"github.com/8h9x/vinderman"
	"github.com/8h9x/vinderman/auth"
	"github.com/8h9x/vinderman/consts"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"gitlab.com/8h9x/Vinderman/eos"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var login = discord.SlashCommandCreate{
	Name:        "login",
	Description: "Login and grant Blast! access to your epic account. (alias for /accounts add)",
}

var Login = Command{
	Handler: func(event *handler.CommandEvent, user database.UserEntry, credentials auth.TokenResponse, data discord.SlashCommandInteractionData) error {
		if len(user.Accounts) >= 25 {
			embed := discord.NewEmbedBuilder().
				SetColor(0xFB5A32).
				SetTimestamp(time.Now()).
				SetAuthorIcon(*event.User().AvatarURL(discord.WithFormat(discord.FileFormatPNG))).
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

		clientCredentials, err := client.Epic.EOS.GetClientCredentials(consts.UEFN_CLIENT_ID, consts.UEFN_CLIENT_SECRET)
		if err != nil {
			return err
		}

		deviceAuthorization, err := client.Epic.EOS.GetDeviceCode(clientCredentials)
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

		deviceCodeCredentials, err := WaitForDeviceCodeAcceptEOS(client.Epic.EOS, consts.UEFN_CLIENT_ID, consts.UEFN_CLIENT_SECRET, deviceAuthorization.DeviceCode, 0)
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

		exchangeCode, err := client.Epic.EOS.GetExchangeCode(deviceCodeCredentials)
		if err != nil {
			return err
		}

		// decode this refresh jwt and get jti for hex format refresh token that expires after ~6 months (170 days) and can be infinitly refreshed
		exchangeCredentials, err := client.Epic.EOS.ExchangeCodeLogin(consts.FORTNITE_PC_CLIENT_ID, consts.FORTNITE_PC_CLIENT_SECRET, exchangeCode.Code)
		if err != nil {
			return err
		}

		decodedRefreshJwtPayload, err := common.Base64Decode(strings.Split(exchangeCredentials.RefreshToken, ".")[1])
		if err != nil {
			return err
		}

		refreshPayload := DecodedRefreshTokenJwtPayload{}

		err = json.Unmarshal(decodedRefreshJwtPayload, &refreshPayload)
		if err != nil {
			return err
		}

		userId := event.User().ID

		_, err = db.Fetch[db.UserEntry]("users", bson.M{"discordId": userId})
		if err != nil {
			return err
		}

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

		refreshCredentials, err := client.Epic.RefreshTokenLogin(consts.FORTNITE_PC_CLIENT_ID, consts.FORTNITE_PC_CLIENT_SECRET, refreshPayload.Jti)
		if err != nil {
			return err
		}

		accountInfo, err := client.Epic.FetchMe(refreshCredentials)
		if err != nil {
			return err
		}

		avatarURL, err := client.Epic.FetchAvatarURL(refreshCredentials)
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

		return err
	},
	LoginRequired:     false,
	EphemeralResponse: false,
}

type DecodedRefreshTokenJwtPayload struct {
	Aud   string `json:"aud"`
	Sub   string `json:"sub"`
	T     string `json:"t"`
	AppID string `json:"appid"`
	Scope string `json:"scope"`
	Iss   string `json:"iss"`
	Dn    string `json:"dn"`
	Exp   int    `json:"exp"`
	Iat   int    `json:"iat"`
	Jti   string `json:"jti"`
	Pfpid string `json:"pfpid"`
}

// this is re-done here because goroutines are so annoying to work with when trying to cancel early
func WaitForDeviceCodeAcceptEOS(c *eos.Client, clientId string, clientSecret string, deviceCode string, runs int) (credentials eos.UserCredentials, err error) {
	if runs == 12 { // after 120 seconds
		return eos.UserCredentials{}, errors.New("errors.blast.timed_out")
	}

	credentials, err = c.DeviceCodeLogin(clientId, clientSecret, deviceCode)

	if err != nil {
		if err.(*request.Error[eos.EpicErrorResponse]).Raw.ErrorCode == "errors.com.epicgames.account.oauth.authorization_pending" {
			time.Sleep(10 * time.Second)
			return WaitForDeviceCodeAcceptEOS(c, clientId, clientSecret, deviceCode, runs+1)
		}

		return
	}

	return
}
