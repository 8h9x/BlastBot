package login

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/8h9x/BlastBot/internal/database"
	"github.com/8h9x/BlastBot/internal/sessions"
	"github.com/8h9x/fortgo/account"
	"github.com/8h9x/fortgo/consts"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"gitlab.com/8h9x/Vinderman/eos"
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
	eosClient := &eos.Client{HttpClient: httpClient}

	//	clientCredentials, err := auth.Authenticate(httpClient, auth.FortnitePS4USClient, auth.PayloadClientCredentials{}, false)
	clientCredentials, err := eosClient.GetClientCredentials(consts.UEFNClientID, consts.UEFNClientSecret)
	if err != nil {
		return err
	}

	//	deviceAuthorization, err := auth.GetDeviceCode(httpClient, clientCredentials)
	deviceAuthorization, err := eosClient.GetDeviceCode(clientCredentials)
	if err != nil {
		return err
	}

	err = event.CreateMessage(discord.MessageCreate{
		Content: fmt.Sprintf("Visit this URL: %s\nThen press 'Confirm' on the epic games login page.\n", deviceAuthorization.VerificationUriComplete),
	})
	if err != nil {
		return err
	}

	deviceCodeCredentials, err := waitForDeviceCodeConfirmEOS(eosClient, consts.UEFNClientID, consts.UEFNClientSecret, deviceAuthorization.DeviceCode, CHECK_INTERVAL, CHECK_TIMEOUT)
	if err != nil {
		return err
	}

	exchangeCode, err := eosClient.GetExchangeCode(deviceCodeCredentials)
	if err != nil {
		return err
	}

	exchangeCredentials, err := eosClient.ExchangeCodeLogin(consts.FortnitePS4USClientID, consts.FortnitePS4USClientSecret, exchangeCode.Code)
	if err != nil {
		return err
	}

	decodedRefreshJWTPayload, err := base64.RawURLEncoding.DecodeString(strings.Split(exchangeCredentials.RefreshToken, ".")[1])
	if err != nil {
		return err
	}

	refreshPayload := DecodedRefreshTokenJwtPayload{}

	err = json.Unmarshal(decodedRefreshJWTPayload, &refreshPayload)
	if err != nil {
		return err
	}

	userId := event.User().ID

	userEntry := database.User{
		ID:        bson.NewObjectID().Hex(),
		DiscordID: userId,
		Accounts: []database.EpicAccount{
			{
				AccountID:        exchangeCredentials.AccountID,
				RefreshToken:     refreshPayload.JTI,
				RefreshExpiresAt: exchangeCredentials.RefreshExpiresAt,
				CreatedClientID:  consts.FortnitePCClientID,
				Flags:            database.USER,
			},
		},
		SelectedEpicAccountId: exchangeCredentials.AccountID,
		BulkFlags:             database.USER,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}

	col := database.GetCollection("users")

	_, err = database.Fetch[database.User]("users", bson.M{"discordId": userId})
	if err == nil { // user exists
		_, err := col.UpdateOne(context.Background(), bson.M{"discordId": userId}, bson.M{"$push": bson.M{"accounts": bson.M{
			"accountId":        exchangeCredentials.AccountID,
			"refreshToken":     refreshPayload.JTI,
			"refreshExpiresAt": exchangeCredentials.RefreshExpiresAt,
			"clientId":         consts.FortnitePCClientID,
			"flags":            database.USER,
		}}}, options.UpdateOne().SetUpsert(true))
		if err != nil {
			return err
		}

		_, err = col.UpdateOne(context.Background(), bson.M{"discordId": userId}, bson.M{"$set": bson.M{"selectedEpicAccountId": exchangeCredentials.AccountID}})
		if err != nil {
			return err
		}
	} else {
		_, err = col.InsertOne(context.Background(), userEntry)
		if err != nil {
			return err
		}
	}

	refreshCredentials, err := account.Authenticate(httpClient, &account.AuthClient{consts.FortnitePS4USClientID, consts.FortnitePS4USClientSecret}, account.TokenPayloadRefreshToken{refreshPayload.JTI}, false)
	if err != nil {
		return err
	}

	session, err := sessions.CreateSession(httpClient, refreshCredentials)
	if err != nil {
		return err
	}

	accountInfo, err := session.AccountService.GetUserByID(session.CurrentCredentials().AccountID)
	if err != nil {
		return err
	}

	avatar, err := session.AvatarService.GetOne(accountInfo.ID)
	if err != nil {
		return err
	}

	avatarURL := fmt.Sprintf("https://fortnite-api.com/images/cosmetics/br/%s/icon.png", strings.Replace(avatar.AvatarID, "ATHENACHARACTER:", "", -1))

	_, err = event.CreateFollowupMessage(discord.NewMessageCreate().
		//		ClearComponents().
		AddEmbeds(discord.NewEmbed().
			WithColor(0xFB5A32).
			WithTimestamp(time.Now()).
			WithAuthorIcon(avatarURL).
			WithAuthorNamef("New account saved for %s", event.User().Username).
			WithDescriptionf("Successfully logged into **%s**\nYou now have (%d/25) saved accounts.", accountInfo.DisplayName, len(userEntry.Accounts)+1),
		),
	)

	return err
}

func waitForDeviceCodeConfirmEOS(c *eos.Client, clientID string, clientSecret string, deviceCode string, interval, timeout time.Duration) (eos.UserCredentials, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return eos.UserCredentials{}, ctx.Err()
		case <-ticker.C:
			//			credentials, err := auth.Authenticate(httpClient, auth.FortnitePS4USClient, payload, true)
			credentials, err := c.DeviceCodeLogin(clientID, clientSecret, deviceCode)
			if err != nil {
				fmt.Println(err)
				continue
			}
			return credentials, nil
		}
	}
}

func waitForDeviceCodeConfirm(httpClient *http.Client, deviceCode string, interval, timeout time.Duration) (account.TokenResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	payload := account.TokenPayloadDeviceCode{
		DeviceCode: deviceCode,
	}

	for {
		select {
		case <-ctx.Done():
			return account.TokenResponse{}, ctx.Err()
		case <-ticker.C:
			credentials, err := account.Authenticate(httpClient, account.FortnitePS4USClient, payload, true)
			if err != nil {
				continue
			}
			return credentials, nil
		}
	}
}

type DecodedRefreshTokenJwtPayload struct {
	AUD   string `json:"aud"`
	SUB   string `json:"sub"`
	T     string `json:"t"`
	AppID string `json:"appid"`
	Scope string `json:"scope"`
	ISS   string `json:"iss"`
	DN    string `json:"dn"`
	EXP   int    `json:"exp"`
	IAT   int    `json:"iat"`
	JTI   string `json:"jti"`
	Pfpid string `json:"pfpid"`
}
