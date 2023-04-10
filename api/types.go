package api

import "time"

type EpicErrorResponse struct {
	ErrorCode          string `json:"errorCode"`
	ErrorMessage       string `json:"errorMessage"`
	NumericErrorCode   int    `json:"numericErrorCode"`
	OriginatingService string `json:"originatingService"`
	Intent             string `json:"intent"`
}

type ClientCredentialsResponse struct {
	AccessToken   string    `json:"access_token"`
	ApplicationId string    `json:"application_id"`
	ClientId      string    `json:"client_id"`
	ExpiresAt     time.Time `json:"expires_at"`
	ExpiresIn     int       `json:"expires_in"`
	TokenType     string    `json:"token_type"`
}

type DeviceAuthorizationResponse struct {
	ClientId                string `json:"client_id"`
	DeviceCode              string `json:"device_code"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`
	Prompt                  string `json:"prompt"`
	UserCode                string `json:"user_code"`
	VerificationUri         string `json:"verification_uri"`
	VerificationUriComplete string `json:"verification_uri_complete"`
}

type UserCredentialsResponse struct {
	AccessToken      string    `json:"access_token"`
	AccountId        string    `json:"account_id"`
	ApplicationId    string    `json:"application_id"`
	ClientId         string    `json:"client_id"`
	ExpiresAt        time.Time `json:"expires_at"`
	ExpiresIn        int       `json:"expires_in"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
	RefreshExpiresIn int       `json:"refresh_expires_in"`
	RefreshToken     string    `json:"refresh_token"`
	Scope            string    `json:"scope"`
	TokenType        string    `json:"token_type"`
}

type ExchangeCodeResponse struct {
	Code             string `json:"code"`
	CreatingClientId string `json:"creatingClientId"`
	ExpiresInSeconds int    `json:"expiresInSeconds"`
}

type DecodedRefreshTokenJwtPayload struct {
	Aud   string `json:"aud"`
	Sub   string `json:"sub"`
	T     string `json:"t"`
	AppID string `json:"appid"`
	Scope string `json:"scope"`
	Iss   string `json:"iss"`
	Dn    string `json:"dn"`
	Exp   int    `json:"exp"`
	Iat   int    `json:"iat"`
	Jti   string `json:"jti"`
	Pfpid string `json:"pfpid"`
}
