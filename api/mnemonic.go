package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c EpicClient) AddFavoriteMnemonic(credentials UserCredentialsResponse, mnemonic string) error {
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Authorization", "Bearer "+credentials.AccessToken)

	resp, err := c.Request("POST", fmt.Sprintf("https://fn-service-discovery-live-public.ogs.live.on.epicgames.com/api/v1/links/favorites/%s/%s", credentials.AccountId, mnemonic), headers, "{}")
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("invalid mnemonic")
	}

	return nil
}

func (c EpicClient) FetchMnemonicInfo(credentials UserCredentialsResponse, mnemonic string) (MnemonicInfoResponse, error) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+credentials.AccessToken)

	resp, err := c.Request("GET", fmt.Sprintf("https://links-public-service-live.ol.epicgames.com/links/api/fn/mnemonic/%s", mnemonic), headers, "")
	if err != nil {
		return MnemonicInfoResponse{}, err
	}

	defer resp.Body.Close()

	var res MnemonicInfoResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return MnemonicInfoResponse{}, err
	}

	return res, nil
}

func (c EpicClient) RemoveFavoriteMnemonic(credentials UserCredentialsResponse, mnemonic string) error {
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Authorization", "Bearer "+credentials.AccessToken)

	resp, err := c.Request("DELETE", fmt.Sprintf("https://fn-service-discovery-live-public.ogs.live.on.epicgames.com/api/v1/links/favorites/%s/%s", credentials.AccountId, mnemonic), headers, "{}")
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("invalid mnemonic")
	}

	return nil
}
