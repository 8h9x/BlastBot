package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c EpicClient) FetchAccountInformation(credentials UserCredentialsResponse) (AccountInformation, error) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+credentials.AccessToken)

	resp, err := c.Request("GET", fmt.Sprintf("https://account-public-service-prod.ol.epicgames.com/account/api/public/account/%s", credentials.AccountId), headers, "")
	if err != nil {
		return AccountInformation{}, err
	}

	defer resp.Body.Close()

	var res AccountInformation
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return AccountInformation{}, err
	}

	return res, nil
}
