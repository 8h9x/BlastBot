package login

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/8h9x/vinderman"
	"github.com/8h9x/vinderman/auth"
	"github.com/8h9x/vinderman/consts"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
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

	clientCredentials, err := auth.Authenticate(httpClient, consts.FortnitePS4USClientID, consts.FortnitePS4USClientSecret, auth.PayloadClientCredentials{}, false)
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

	_, err = vinderman.NewClient(httpClient, credentials)
	if err != nil {
		return err
	}

	_, err = event.CreateFollowupMessage(discord.MessageCreate{
		Content: "Vinderman client successfully created",
	})

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
			credentials, err := auth.Authenticate(httpClient, consts.FortnitePS4USClientID, consts.FortnitePS4USClientSecret, payload, true)
			if err != nil {
				continue
			}
			return credentials, nil
		}
	}
}
