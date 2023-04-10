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
