package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type EpicClient struct {
	httpClient *http.Client
}

func New() *EpicClient {
	return &EpicClient{
		httpClient: &http.Client{},
	}
}

type RequestError struct {
	StatusCode int
	Message    string
	Raw        EpicErrorResponse
}

func (e *RequestError) Error() string {
	return e.Message
}

func (c EpicClient) Request(method string, url string, headers http.Header, body string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "FortniteGame/++Fortnite+Release-24.10-CL-24900093 Windows/10.0.22000.1.768.64bit")

	for key, value := range headers {
		req.Header.Set(key, value[0])
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		if resp.Body != nil {
			defer resp.Body.Close()

			var res EpicErrorResponse
			err = json.NewDecoder(resp.Body).Decode(&res)
			if err != nil {
				return nil, err
			}

			if res.ErrorMessage != "" {
				return nil, &RequestError{
					StatusCode: resp.StatusCode,
					Message:    fmt.Sprintf("%s request to %s failed with error message: %s", method, url, res.ErrorMessage),
					Raw:        res,
				}
			}
		}

		return nil, fmt.Errorf("%s request to %s failed with status code %d", method, url, resp.StatusCode)
	}

	return resp, nil
}

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

func (c EpicClient) DeviceCodeLoginEOS(clientId string, clientSecret string, deviceCode string) (UserCredentialsResponse, error) {
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

func (c EpicClient) WaitForDeviceCodeAcceptEOS(clientId string, clientSecret string, deviceCode string) (UserCredentialsResponse, error) {
	deviceCodeCredentials, err := c.DeviceCodeLoginEOS(clientId, clientSecret, deviceCode)

	if err != nil {
		if err.(*RequestError).Raw.ErrorCode == "errors.com.epicgames.account.oauth.authorization_pending" {
			time.Sleep(3 * time.Second)
			return c.WaitForDeviceCodeAcceptEOS(clientId, clientSecret, deviceCode)
		}

		return UserCredentialsResponse{}, err
	}

	return deviceCodeCredentials, nil
}

// https://account-public-service-prod.ol.epicgames.com/account/api/oauth/exchange
func (c EpicClient) GetExchangeCodeEOS(credentials UserCredentialsResponse) (ExchangeCodeResponse, error) {
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

func (c EpicClient) ExchangeCodeLoginEOS(clientId string, clientSecret string, exchangeCode string) (UserCredentialsResponse, error) {
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
