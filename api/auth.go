package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// TODO: Massive amount of de-duplication needed here

func (c EpicClient) GetClientCredentialsEOS(clientId string, clientSecret string) (ClientCredentialsResponse, error) {
	encodedClientToken := base64Encode(clientId + ":" + clientSecret)

	headers := http.Header{}
	headers.Set("Content-Type", "application/x-www-form-urlencoded")
	headers.Set("Authorization", fmt.Sprint("Basic ", encodedClientToken))

	v := url.Values{}
	v.Set("grant_type", "client_credentials")
	body := v.Encode()

	resp, err := c.Request("POST", "https://api.epicgames.dev/epic/oauth/v2/token", headers, body)
	if err != nil {
		return ClientCredentialsResponse{}, err
	}

	defer resp.Body.Close()

	var res ClientCredentialsResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return ClientCredentialsResponse{}, err
	}

	return res, nil
}

func (c EpicClient) GetDeviceAuthorizationEOS(r string) (DeviceAuthorizationResponse, error) {
	headers := http.Header{}
	headers.Set("Content-Type", "application/x-www-form-urlencoded")
	headers.Set("Authorization", fmt.Sprint("Bearer ", r))

	v := url.Values{}
	v.Set("prompt", "login")
	body := v.Encode()

	resp, err := c.Request("POST", "https://api.epicgames.dev/epic/oauth/v2/deviceAuthorization", headers, body)
	if err != nil {
		return DeviceAuthorizationResponse{}, err
	}

	defer resp.Body.Close()

	var res DeviceAuthorizationResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return DeviceAuthorizationResponse{}, err
	}

	return res, nil
}

func (c EpicClient) DeviceCodeLoginEOS(clientId string, clientSecret string, deviceCode string) (UserCredentialsResponseEOS, error) {
	encodedClientToken := base64Encode(clientId + ":" + clientSecret)

	headers := http.Header{}
	headers.Set("Content-Type", "application/x-www-form-urlencoded")
	headers.Set("Authorization", fmt.Sprint("Basic ", encodedClientToken))

	v := url.Values{}
	v.Set("grant_type", "device_code")
	v.Set("device_code", deviceCode)
	body := v.Encode()

	resp, err := c.Request("POST", "https://api.epicgames.dev/epic/oauth/v2/token", headers, body)
	if err != nil {
		return UserCredentialsResponseEOS{}, err
	}

	defer resp.Body.Close()

	var res UserCredentialsResponseEOS
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return UserCredentialsResponseEOS{}, err
	}

	return res, nil
}

func (c EpicClient) WaitForDeviceCodeAcceptEOS(clientId string, clientSecret string, deviceCode string) (UserCredentialsResponseEOS, error) {
	deviceCodeCredentials, err := c.DeviceCodeLoginEOS(clientId, clientSecret, deviceCode)

	if err != nil {
		if err.(*RequestError).Raw.ErrorCode == "errors.com.epicgames.account.oauth.authorization_pending" {
			time.Sleep(10 * time.Second)
			return c.WaitForDeviceCodeAcceptEOS(clientId, clientSecret, deviceCode)
		}

		return UserCredentialsResponseEOS{}, err
	}

	return deviceCodeCredentials, nil
}

// https://account-public-service-prod.ol.epicgames.com/account/api/oauth/exchange
func (c EpicClient) GetExchangeCodeEOS(credentials UserCredentialsResponseEOS) (ExchangeCodeResponse, error) {
	headers := http.Header{}
	headers.Set("Authorization", fmt.Sprint("Bearer ", credentials.AccessToken))

	resp, err := c.Request("GET", "https://api.epicgames.dev/epic/oauth/v2/exchange", headers, "")
	if err != nil {
		return ExchangeCodeResponse{}, err
	}

	defer resp.Body.Close()

	var res ExchangeCodeResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return ExchangeCodeResponse{}, err
	}

	return res, nil
}

func (c EpicClient) ExchangeCodeLoginEOS(clientId string, clientSecret string, exchangeCode string) (UserCredentialsResponseEOS, error) {
	encodedClientToken := base64Encode(clientId + ":" + clientSecret)

	headers := http.Header{}
	headers.Set("Content-Type", "application/x-www-form-urlencoded")
	headers.Set("Authorization", fmt.Sprint("Basic ", encodedClientToken))

	v := url.Values{}
	v.Set("grant_type", "exchange_code")
	v.Set("exchange_code", exchangeCode)
	v.Set("scope", "offline_access")
	body := v.Encode()

	resp, err := c.Request("POST", "https://api.epicgames.dev/epic/oauth/v2/token", headers, body)
	if err != nil {
		return UserCredentialsResponseEOS{}, err
	}

	defer resp.Body.Close()

	var res UserCredentialsResponseEOS
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return UserCredentialsResponseEOS{}, err
	}

	return res, nil
}

func (c EpicClient) RefreshTokenLogin(clientId string, clientSecret string, refreshToken string) (UserCredentialsResponse, error) {
	encodedClientToken := base64Encode(clientId + ":" + clientSecret)

	headers := http.Header{}
	headers.Set("Content-Type", "application/x-www-form-urlencoded")
	headers.Set("Authorization", fmt.Sprint("Basic ", encodedClientToken))

	v := url.Values{}
	v.Set("grant_type", "refresh_token")
	v.Set("refresh_token", refreshToken)
	body := v.Encode()

	resp, err := c.Request("POST", "https://account-public-service-prod.ol.epicgames.com/account/api/oauth/token", headers, body)
	if err != nil {
		return UserCredentialsResponse{}, err
	}

	defer resp.Body.Close()

	var res UserCredentialsResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return UserCredentialsResponse{}, err
	}

	return res, nil
}

func (c EpicClient) GetClientCredentials(clientId string, clientSecret string) (ClientCredentialsResponse, error) {
	encodedClientToken := base64Encode(clientId + ":" + clientSecret)

	headers := http.Header{}
	headers.Set("Content-Type", "application/x-www-form-urlencoded")
	headers.Set("Authorization", fmt.Sprint("Basic ", encodedClientToken))

	v := url.Values{}
	v.Set("grant_type", "client_credentials")
	body := v.Encode()

	resp, err := c.Request("POST", "https://account-public-service-prod.ol.epicgames.com/account/api/oauth/token", headers, body)
	if err != nil {
		return ClientCredentialsResponse{}, err
	}

	defer resp.Body.Close()

	var res ClientCredentialsResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return ClientCredentialsResponse{}, err
	}

	return res, nil
}

func (c EpicClient) GetExchangeCode(credentials UserCredentialsResponse) (ExchangeResponse, error) {
	headers := http.Header{}
	headers.Set("Authorization", fmt.Sprint("Bearer ", credentials.AccessToken))

	resp, err := c.Request("GET", "https://account-public-service-prod.ol.epicgames.com/account/api/oauth/exchange", headers, "")
	if err != nil {
		return ExchangeResponse{}, err
	}

	defer resp.Body.Close()

	var res ExchangeResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return ExchangeResponse{}, err
	}

	return res, nil
}

func (c EpicClient) GetDeviceAuth(credentials UserCredentialsResponse) (DeviceAuth, error) {
	headers := http.Header{}
	headers.Set("Authorization", fmt.Sprint("Bearer ", credentials.AccessToken))

	resp, err := c.Request("POST", fmt.Sprintf("https://account-public-service-prod.ol.epicgames.com/account/api/public/account/%s/deviceAuth", credentials.AccountID), headers, "")
	if err != nil {
		return DeviceAuth{}, err
	}

	defer resp.Body.Close()

	var res DeviceAuth
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return DeviceAuth{}, err
	}

	return res, nil
}

func (c EpicClient) ExchangeCodeLogin(clientId string, clientSecret string, exchangeCode string) (UserCredentialsResponse, error) {
	encodedClientToken := base64Encode(clientId + ":" + clientSecret)

	headers := http.Header{}
	headers.Set("Content-Type", "application/x-www-form-urlencoded")
	headers.Set("Authorization", fmt.Sprint("Basic ", encodedClientToken))

	v := url.Values{}
	v.Set("grant_type", "exchange_code")
	v.Set("exchange_code", exchangeCode)
	v.Set("scope", "offline_access")
	body := v.Encode()

	resp, err := c.Request("POST", "https://account-public-service-prod.ol.epicgames.com/account/api/oauth/token", headers, body)
	if err != nil {
		return UserCredentialsResponse{}, err
	}

	defer resp.Body.Close()

	var res UserCredentialsResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return UserCredentialsResponse{}, err
	}

	return res, nil
}

func base64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}
