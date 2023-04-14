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

type UserCredentialsResponseEOS struct {
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
	Scope            []string  `json:"scope"`
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

type MnemonicInfoResponse struct {
	AccountID       string    `json:"accountId"`
	Active          bool      `json:"active"`
	Created         time.Time `json:"createdAt"`
	CreatorName     string    `json:"creatorName"`
	DescriptionTags []string  `json:"descriptionTags"`
	Disabled        bool      `json:"disabled"`
	LinkType        string    `json:"linkType"`
	Metadata        struct {
		DynamicXp struct {
			CalibrationPhase  string `json:"calibrationPhase"`
			UniqueGameVersion int    `json:"uniqueGameVersion"`
		} `json:"dynamicXp"`
		GeneratedIslandUrlsOld struct {
			URL  string `json:"url"`
			URLM string `json:"url_m"`
			URLS string `json:"url_s"`
		} `json:"generated_island_urls_old"`
		Introduction string `json:"introduction"`
		IslandType   string `json:"islandType"`
		Locale       string `json:"locale"`
		Matchmaking  struct {
			BAllowJoinInProgress   bool   `json:"bAllowJoinInProgress"`
			JoinInProgressTeam     int    `json:"joinInProgressTeam"`
			JoinInProgressType     string `json:"joinInProgressType"`
			MaximumNumberOfPlayers int    `json:"maximumNumberOfPlayers"`
			MinimumNumberOfPlayers int    `json:"minimumNumberOfPlayers"`
			MmsType                string `json:"mmsType"`
			NumberOfTeams          int    `json:"numberOfTeams"`
			OverridePlaylist       string `json:"override_Playlist"`
			PlayerCount            int    `json:"playerCount"`
			PlayersPerTeam         int    `json:"playersPerTeam"`
		} `json:"matchmaking"`
		SupportCode string `json:"supportCode"`
		Tagline     string `json:"tagline"`
		Title       string `json:"title"`
	} `json:"metadata"`
	Mnemonic         string `json:"mnemonic"`
	ModerationStatus string `json:"moderationStatus"`
	Namespace        string `json:"namespace"`
	Version          int    `json:"version"`
}

type AccountInformation struct {
	AgeGroup                   string    `json:"ageGroup"`
	CabinedMode                bool      `json:"cabinedMode"`
	CanUpdateDisplayName       bool      `json:"canUpdateDisplayName"`
	CanUpdateDisplayNameNext   time.Time `json:"canUpdateDisplayNameNext"`
	Country                    string    `json:"country"`
	DisplayName                string    `json:"displayName"`
	Email                      string    `json:"email"`
	EmailVerified              bool      `json:"emailVerified"`
	FailedLoginAttempts        int       `json:"failedLoginAttempts"`
	HasHashedEmail             bool      `json:"hasHashedEmail"`
	Headless                   bool      `json:"headless"`
	Id                         string    `json:"id"`
	LastDisplayNameChange      time.Time `json:"lastDisplayNameChange"`
	LastLogin                  time.Time `json:"lastLogin"`
	LastName                   string    `json:"lastName"`
	MinorExpected              bool      `json:"minorExpected"`
	MinorStatus                string    `json:"minorStatus"`
	MinorVerified              bool      `json:"minorVerified"`
	Name                       string    `json:"name"`
	NumberOfDisplayNameChanges int       `json:"numberOfDisplayNameChanges"`
	PreferredLanguage          string    `json:"preferredLanguage"`
	TfaEnabled                 bool      `json:"tfaEnabled"`
}

type ExchangeResponse struct {
	Code             string `json:"code"`
	CreatingClientId string `json:"creatingClientId"`
	ExpiresInSeconds int    `json:"expiresInSeconds"`
}

type DeviceAuth struct {
	AccountId string `json:"accountId"`
	Created   struct {
		DateTime  time.Time `json:"dateTime"`
		IpAddress string    `json:"ipAddress"`
		Location  string    `json:"location"`
	} `json:"created"`
	DeviceId  string `json:"deviceId"`
	Secret    string `json:"secret"`
	UserAgent string `json:"userAgent"`
}

type BRInfo struct {
	Stash struct {
		Globalcash int `json:"globalcash"`
	} `json:"stash"`
}

/*

map[
	accountId:3de654385bf64a9686614cdc999d2d21
	created:map[
		dateTime:2023-04-14T01:57:33.714Z
		ipAddress:65.190.218.12
		location:Bladenboro,United States
	]
	deviceId:b111059de97a465d93fee59fabad7c86
	secret:WNAG5XCF756MMPGFDOWTYHU3J6E65YHZ
	userAgent:FortniteGame/++Fortnite+Release-24.10-CL-24900093 Windows/10.0.22000.1.768.64bit
]

*/

/*
map[
	ageGroup:UNKNOWN
	cabinedMode:false
	canUpdateDisplayName:false
	canUpdateDisplayNameNext:2023-04-27T01:15:10.851Z
	country:US
	displayName:NarcissisticApe
	email:tom.haze420@yahoo.com
	emailVerified:true
	failedLoginAttempts:0
	hasHashedEmail:false
	headless:false
	id:3de654385bf64a9686614cdc999d2d21
	lastDisplayNameChange:2023-04-13T01:15:10.851Z
	lastLogin:2023-04-11T02:24:23.330Z
	lastName:Burggess
	minorExpected:false
	minorStatus:NOT_MINOR
	minorVerified:false
	name:Sabrina
	numberOfDisplayNameChanges:7
	preferredLanguage:en
	tfaEnabled:true
]
*/

/*

map[
	accountId:e0121ce665474edfab586ce98ab791ed
	active:true
	created:2019-03-29T20:05:34.195Z
	creatorName:brendannnd
	descriptionTags:[1v1 free for all pvp]
	disabled:false
	linkType:Creative:Island
	metadata:map[
		dynamicXp:map[
			calibrationPhase:DataGathering
			uniqueGameVersion:1
		]
		generated_island_urls_old:map[
			url:https://fortnite-island-screenshots-live-cdn.ol.epicgames.com/screenshots/5224-2376-0131.png
			url_m:https://fortnite-island-screenshots-live-cdn.ol.epicgames.com/screenshots/5224-2376-0131_m.png
			url_s:https://fortnite-island-screenshots-live-cdn.ol.epicgames.com/screenshots/5224-2376-0131_s.png
		]
		introduction:Endless 1v1 build fights with any guns!
		islandType:CreativePlot:temperate_medium
		locale:en
		matchmaking:map[
			bAllowJoinInProgress:true
			joinInProgressTeam:255
			joinInProgressType:JoinImmediately
			maximumNumberOfPlayers:10
			minimumNumberOfPlayers:10
			mmsType:keep_full
			numberOfTeams:10
			override_Playlist:
			playerCount:10
			playersPerTeam:-1
		]
		supportCode:brendannnd
		tagline:Endless 1v1 build fights with any guns!
		title:BRENDAN'S 1v1 BUILD FIGHTS
	]
	mnemonic:1111-1111-1111
	moderationStatus:Approved
	namespace:fn
	version:1
]

*/
