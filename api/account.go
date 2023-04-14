package api

import (
	"blast/api/consts"
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

func (c EpicClient) FetchAccountBRInfo(credentials UserCredentialsResponse) (BRInfo, error) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+credentials.AccessToken)

	resp, err := c.Request("GET", fmt.Sprintf("%s/br-inventory/account/%s", consts.FORTNITE_GAME_BASE, credentials.AccountId), headers, "")
	if err != nil {
		return BRInfo{}, err
	}

	defer resp.Body.Close()

	var res BRInfo
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return BRInfo{}, err
	}

	return res, nil
}

func (c EpicClient) FetchPartyInfo(credentials UserCredentialsResponse) (any, error) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+credentials.AccessToken)

	resp, err := c.Request("GET", fmt.Sprintf("https://party-service-prod.ol.epicgames.com/party/api/v1/Fortnite/user/%s", credentials.AccountId), headers, "")
	if err != nil {
		return BRInfo{}, err
	}

	defer resp.Body.Close()

	var res any
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return BRInfo{}, err
	}

	return res, nil
}

// https://party-service-prod.ol.epicgames.com/party/api/v1/Fortnite/user/accountid
