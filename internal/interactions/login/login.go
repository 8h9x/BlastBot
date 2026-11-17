package login

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/8h9x/BlastBot/internal/database"
	"github.com/8h9x/BlastBot/internal/sessions"
	"github.com/8h9x/fortgo"
	"github.com/8h9x/fortgo/account"
	"github.com/8h9x/fortgo/consts"
	"github.com/8h9x/fortgo/eos"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	CHECK_INTERVAL = time.Second * 10
	CHECK_TIMEOUT  = time.Minute * 2
)

const (
	interval = time.Second * 10
	timeout  = time.Minute * 2
)

var Definition = discord.SlashCommandCreate{
	Name:        "login",
	Description: "Login and grant Blast! access to your epic account. (alias for /accounts add)",
}

type TemporaryTokenSource struct {
	token eos.Token
}

func (s *TemporaryTokenSource) Token(ctx context.Context) (string, error) {
	return s.token.AccessToken, nil
}

func Handler(event *handler.CommandEvent) error {
	src := &TemporaryTokenSource{}
	client := eos.New(fortgo.NewHTTPClient(src))

	clientToken, err := client.GetToken(
		&account.AuthClient{
			ID:     consts.UEFNClientID,
			Secret: consts.UEFNClientSecret,
		},
		account.TokenGrantClientCredentials(),
		account.TokenTypeDefault,
	)
	if err != nil {
		return err
	}

	deviceAuthorization, err := client.CreateDeviceCode(clientToken)
	if err != nil {
		return err
	}

	err = event.CreateMessage(discord.MessageCreate{
		Content: fmt.Sprintf("Visit this URL: %s\nThen press 'Confirm' on the epic games login page.\n", deviceAuthorization.VerificationURIComplete),
	})
	if err != nil {
		return err
	}

	src.token, err = waitForConfirmAndGetToken(
		client,
		&account.AuthClient{
			ID:     consts.UEFNClientID,
			Secret: consts.UEFNClientSecret,
		},
		deviceAuthorization.DeviceCode,
		interval,
		timeout,
	)
	if err != nil {
		return err
	}

	ex, err := client.GetExchangeCode()
	if err != nil {
		return err
	}

	token, err := client.GetToken(
		account.FortnitePS4USClient,
		account.TokenGrantExchangeCode(ex.Code),
		account.TokenTypeDefault,
		"offline_access",
	)
	if err != nil {
		return err
	}

	decodedRefreshJWTPayload, err := base64.RawURLEncoding.DecodeString(strings.Split(token.RefreshToken, ".")[1])
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
				AccountID:        token.AccountID,
				RefreshToken:     refreshPayload.JTI,
				RefreshExpiresAt: token.RefreshExpiresAt,
				CreatedClientID:  consts.FortnitePCClientID,
				Flags:            database.USER,
			},
		},
		SelectedEpicAccountId: token.AccountID,
		BulkFlags:             database.USER,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}

	col := database.GetCollection("users")

	_, err = database.Fetch[database.User]("users", bson.M{"discordId": userId})
	if err == nil { // user exists
		_, err := col.UpdateOne(context.Background(), bson.M{"discordId": userId}, bson.M{"$push": bson.M{"accounts": bson.M{
			"accountId":        token.AccountID,
			"refreshToken":     refreshPayload.JTI,
			"refreshExpiresAt": token.RefreshExpiresAt,
			"clientId":         consts.FortnitePCClientID,
			"flags":            database.USER,
		}}}, options.UpdateOne().SetUpsert(true))
		if err != nil {
			return err
		}

		_, err = col.UpdateOne(context.Background(), bson.M{"discordId": userId}, bson.M{"$set": bson.M{"selectedEpicAccountId": token.AccountID}})
		if err != nil {
			return err
		}
	} else {
		_, err = col.InsertOne(context.Background(), userEntry)
		if err != nil {
			return err
		}
	}

	session, err := sessions.CreateSession(account.TokenGrantRefreshToken(refreshPayload.JTI))
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

func waitForConfirmAndGetToken(client *eos.Client, authClient *account.AuthClient, deviceCode string, interval, timeout time.Duration) (eos.Token, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return eos.Token{}, ctx.Err()
		case <-ticker.C:
			token, err := client.GetToken(
				&account.AuthClient{
					ID:     consts.UEFNClientID,
					Secret: consts.UEFNClientSecret,
				},
				account.TokenGrantClientCredentials(),
				account.TokenTypeDefault,
			)
			if err != nil {
				log.Println(err)
				continue
			}

			return token, nil
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
