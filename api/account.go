package api

import (
	"blast/api/consts"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c EpicClient) FetchAccountInformation(credentials UserCredentialsResponse) (AccountInformation, error) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+credentials.AccessToken)

	resp, err := c.Request("GET", fmt.Sprintf("https://account-public-service-prod.ol.epicgames.com/account/api/public/account/%s", credentials.AccountID), headers, "")
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

	resp, err := c.Request("GET", fmt.Sprintf("%s/br-inventory/account/%s", consts.FORTNITE_GAME_BASE, credentials.AccountID), headers, "")
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

func (c EpicClient) FetchAvatars(credentials UserCredentialsResponse, accountIDs ...string) (FetchAvatarsResponse, error) {
	if len(accountIDs) == 0 {
		accountIDs = append(accountIDs, credentials.AccountID)
	}

	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+credentials.AccessToken)

	resp, err := c.Request("GET", fmt.Sprintf("https://avatar-service-prod.identity.live.on.epicgames.com/v1/avatar/fortnite/ids?accountIds=%s", strings.Join(accountIDs, ",")), headers, "")
	if err != nil {
		return FetchAvatarsResponse{}, err
	}

	defer resp.Body.Close()

	var res FetchAvatarsResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return FetchAvatarsResponse{}, err
	}

	return res, nil
}

func (c EpicClient) FetchAvatar(credentials UserCredentialsResponse) (AccountAvatar, error) {
	avatars, err := c.FetchAvatars(credentials)
	if err != nil {
		return AccountAvatar{}, err
	}

	return avatars[0], nil
}

func (c EpicClient) FetchAvatarURL(credentials UserCredentialsResponse) (string, error) {
	avatar, err := c.FetchAvatar(credentials)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://fortnite-api.com/images/cosmetics/br/%s/icon.png", strings.Replace(avatar.AvatarID, "ATHENACHARACTER:", "", -1)), nil
}
