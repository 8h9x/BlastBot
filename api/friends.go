package api

import (
	"fmt"
	"net/http"
)

func (c EpicClient) AddFriend(credentials UserCredentialsResponse, friendID string) error {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+credentials.AccessToken)

	resp, err := c.Request("POST", fmt.Sprintf("https://friends-public-service-prod.ol.epicgames.com/friends/api/v1/%s/friends/%s", credentials.AccountID, friendID), headers, "")
	if err != nil {
		return err
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("failed to add friend: %s", resp.Status)
	}

	return nil
}
