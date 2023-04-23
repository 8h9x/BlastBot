package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func (c EpicClient) UserSearch(credentials UserCredentialsResponse, displayName string, platform string) (UserSearchResponse, error) {
	values := url.Values{}
	values.Add("prefix", displayName)
	values.Add("platform", platform)

	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+credentials.AccessToken)

	resp, err := c.Request("GET", fmt.Sprintf("https://user-search-service-prod.ol.epicgames.com/api/v1/search/%s?%s", credentials.AccountID, values.Encode()), headers, "")
	if err != nil {
		return UserSearchResponse{}, err
	}

	defer resp.Body.Close()

	var res UserSearchResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return UserSearchResponse{}, err
	}

	return res, nil
}

type UserSearch struct {
	AccountID string `json:"accountId"`
	Matches   []struct {
		Value    string `json:"value"`
		Platform string `json:"platform"`
	} `json:"matches"`
	MatchType    string `json:"matchType"`
	EpicMutuals  int    `json:"epicMutuals"`
	SortPosition int    `json:"sortPosition"`
}

type UserSearchResponse []UserSearch

// https://user-search-service-prod.ol.epicgames.com/api/v1/search/{AccountID}?prefix={displayName}&platform{epic | psn | xbl | steam }

func (c EpicClient) FetchUserByAccountID(credentials UserCredentialsResponse, accountID string) (AccountInformation, error) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+credentials.AccessToken)

	resp, err := c.Request("GET", fmt.Sprintf("https://account-public-service-prod.ol.epicgames.com/account/api/public/account/%s", accountID), headers, "")
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

func (c EpicClient) FetchUserByDisplayName(credentials UserCredentialsResponse, displayName string) (AccountInformation, error) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+credentials.AccessToken)

	resp, err := c.Request("GET", fmt.Sprintf("https://account-public-service-prod.ol.epicgames.com/account/api/public/account/displayName/%s", displayName), headers, "")
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

func (c EpicClient) FetchUserByExternalDisplayName(credentials UserCredentialsResponse, displayName string, platform string) (ExternalDisplayNameLookup, error) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+credentials.AccessToken)

	resp, err := c.Request("GET", fmt.Sprintf("https://account-public-service-prod.ol.epicgames.com/account/api/public/account/lookup/externalAuth/%s/displayName/%s?caseInsensitive=true", platform, displayName), headers, "")
	if err != nil {
		return ExternalDisplayNameLookup{}, err
	}

	defer resp.Body.Close()

	var res ExternalDisplayNameLookup
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return ExternalDisplayNameLookup{}, err
	}

	return res, nil
}

type ExternalDisplayNameLookup []struct {
	ID            string `json:"id"`
	DisplayName   string `json:"displayName"`
	ExternalAuths map[string]struct {
		AccountID           string        `json:"accountId"`
		Type                string        `json:"type"`
		ExternalAuthIDType  string        `json:"externalAuthIdType"`
		ExternalDisplayName string        `json:"externalDisplayName"`
		AuthIds             []interface{} `json:"authIds"`
	} `json:"externalAuths"`
}

// https://github.com/LeleDerGrasshalmi/FortniteEndpointsDocumentation/blob/d754c550ad762576a17fc5684f0fc7a267d9a5e7/EpicGames/AccountService/Account/ExternalAuth/README.md
