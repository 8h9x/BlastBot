package sessions

import (
	"fmt"
	"net/http"

	"github.com/8h9x/BlastBot/database/internal/database"
	"github.com/8h9x/vinderman"
	"github.com/8h9x/vinderman/auth"
	"github.com/8h9x/vinderman/consts"
	"github.com/disgoorg/snowflake/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// TODO: At some point it would probably make sense to move some of the auth invalidation logic and coupled event handlers upstream to vinderman
//

var sessions = make(map[string]*vinderman.Client, 0)

func CreateSession(httpClient *http.Client, credentials auth.TokenResponse) (*vinderman.Client, error) {
	session, err := vinderman.NewClient(&http.Client{}, credentials)
	if err != nil {
		return &vinderman.Client{}, fmt.Errorf("an error occured when creating vinderman client %s", err)
	}

	sessions[credentials.AccountID] = session
	return session, nil
}

func GetSession(accountId string) *vinderman.Client {
	return sessions[accountId]
}

func GetSessionForUser(discordId snowflake.ID) (*vinderman.Client, error) {
	user, err := database.Fetch[database.User]("users", bson.M{"discordId": discordId})
	if err != nil {
		return &vinderman.Client{}, fmt.Errorf("unable to fetch user data from database: %s", err)
	}

	_, ok := sessions[user.SelectedEpicAccountId]
	if !ok {
		var activeAccount database.EpicAccount
		for _, account := range user.Accounts {
			if account.AccountID == user.SelectedEpicAccountId {
				activeAccount = account
			}
		}

		httpClient := &http.Client{}

		credentials, err := auth.Authenticate(httpClient, consts.FortnitePS4USClientID, consts.FortnitePS4USClientSecret, auth.PayloadRefreshToken{
			RefreshToken: activeAccount.RefreshToken,
		}, true)
		if err != nil {
			return &vinderman.Client{}, fmt.Errorf("unable to generate credentials using saved refresh token: %s", err)
		}

		CreateSession(httpClient, credentials)
	}

	return sessions[user.SelectedEpicAccountId], nil
}
