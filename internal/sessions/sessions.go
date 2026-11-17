package sessions

import (
	"fmt"
	"net/http"

	"github.com/8h9x/BlastBot/internal/database"
	"github.com/8h9x/fortgo"
	"github.com/8h9x/fortgo/account"
	"github.com/disgoorg/snowflake/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// TODO: At some point it would probably make sense to move some of the auth invalidation logic and coupled event handlers upstream to vinderman
//

var sessions = make(map[string]*fortgo.Client, 0)

func CreateSession(tokenGrant account.TokenGrant) (*fortgo.Client, error) {
	//	session := fortgo.NewClient(&http.Client{}, credentials)
	session := fortgo.NewClient(
		fortgo.WithTokenGrant(tokenGrant),
	)

	err := session.Connect()
	if err != nil {
		return &fortgo.Client{}, fmt.Errorf("an error occured when creating fortgo client %s", err)
	}

	sessions[session.CurrentCredentials().AccountID] = session
	return session, nil
}

func GetSession(accountId string) *fortgo.Client {
	return sessions[accountId]
}

func GetSessionForUser(discordId snowflake.ID) (*fortgo.Client, error) {
	user, err := database.Fetch[database.User]("users", bson.M{"discordId": discordId})
	if err != nil {
		return &fortgo.Client{}, fmt.Errorf("unable to fetch user data from database: %s", err)
	}

	_, ok := sessions[user.SelectedEpicAccountId]
	if !ok {
		var activeAccount database.EpicAccount
		for _, account := range user.Accounts {
			if account.AccountID == user.SelectedEpicAccountId {
				activeAccount = account
			}
		}

		session := fortgo.NewClient(
			fortgo.WithHttpClient(http.DefaultClient),
			fortgo.WithTokenGrant(account.TokenGrantRefreshToken(activeAccount.RefreshToken)),
		)

		err := session.Connect()
		if err != nil {
			return &fortgo.Client{}, fmt.Errorf("an error occured when creating fortgo client %s", err)
		}

		sessions[user.SelectedEpicAccountId] = session
	}

	return sessions[user.SelectedEpicAccountId], nil
}
