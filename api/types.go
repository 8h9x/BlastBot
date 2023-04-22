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
	AccountID        string    `json:"account_id"`
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
	AccountID        string    `json:"account_id"`
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
			CalibrationPhase string `json:"calibrationPhase"`
			// UniqueGameVersion string `json:"uniqueGameVersion"` // somtimes int, sometimes string
		} `json:"dynamicXp"`
		GeneratedIslandUrlsOld struct {
			URL  string `json:"url"`
			URLM string `json:"url_m"`
			URLS string `json:"url_s"`
		} `json:"generated_island_urls_old"`
		ImageURL     string `json:"image_url"`
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
	Mnemonic         string    `json:"mnemonic"`
	ModerationStatus string    `json:"moderationStatus"`
	Namespace        string    `json:"namespace"`
	Published        time.Time `json:"published"`
	Version          int       `json:"version"`
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
	AccountID string `json:"accountId"`
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

type APIVersionResponse struct {
	App                       string    `json:"app"`
	ServerDate                time.Time `json:"serverDate"`
	OverridePropertiesVersion string    `json:"overridePropertiesVersion"`
	CLN                       string    `json:"cln"`
	Build                     string    `json:"build"`
	ModuleName                string    `json:"moduleName"`
	BuildDate                 time.Time `json:"buildDate"`
	Version                   string    `json:"version"`
	Branch                    string    `json:"branch"`
	Modules                   map[string]struct {
		CLN       string    `json:"cln"`
		Build     string    `json:"build"`
		BuildDate time.Time `json:"buildDate"`
		Version   string    `json:"version"`
		Branch    string    `json:"branch"`
	} `json:"modules"`
}

// [{"accountId":"3de654385bf64a9686614cdc999d2d21","namespace":"fortnite","avatarId":"ATHENACHARACTER:CID_083_ATHENA_COMMANDO_F_TACTICAL"}]

type AccountAvatar struct {
	AccountID string `json:"accountId"`
	Namespace string `json:"namespace"`
	AvatarID  string `json:"avatarId"`
}

type FetchAvatarsResponse []AccountAvatar

type OfferResponse map[string]struct {
	ID              string `json:"id"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	LongDescription string `json:"longDescription"`
	KeyImages       []struct {
		Type         string    `json:"type"`
		URL          string    `json:"url"`
		MD5          string    `json:"md5"`
		Width        int       `json:"width"`
		Height       int       `json:"height"`
		Size         int       `json:"size"`
		UploadedDate time.Time `json:"uploadedDate"`
	} `json:"keyImages"`
	Categories []struct {
		Path string `json:"path"`
	} `json:"categories"`
	Namespace        string    `json:"namespace"`
	Status           string    `json:"status"`
	CreationDate     time.Time `json:"creationDate"`
	LastModifiedDate time.Time `json:"lastModifiedDate"`
	CustomAttributes struct {
	} `json:"customAttributes"`
	InternalName string `json:"internalName"`
	Recurrence   string `json:"recurrence"`
	Items        []struct {
		ID              string        `json:"id"`
		Title           string        `json:"title"`
		Description     string        `json:"description"`
		LongDescription string        `json:"longDescription"`
		KeyImages       []interface{} `json:"keyImages"`
		Categories      []struct {
			Path string `json:"path"`
		} `json:"categories"`
		Namespace           string        `json:"namespace"`
		Status              string        `json:"status"`
		CreationDate        time.Time     `json:"creationDate"`
		LastModifiedDate    time.Time     `json:"lastModifiedDate"`
		EntitlementName     string        `json:"entitlementName"`
		EntitlementType     string        `json:"entitlementType"`
		ItemType            string        `json:"itemType"`
		ReleaseInfo         []interface{} `json:"releaseInfo"`
		Developer           string        `json:"developer"`
		DeveloperID         string        `json:"developerId"`
		UseCount            int           `json:"useCount"`
		EulaIds             []interface{} `json:"eulaIds"`
		EndOfSupport        bool          `json:"endOfSupport"`
		NsMajorItems        []interface{} `json:"nsMajorItems"`
		NsDependsOnDlcItems []interface{} `json:"nsDependsOnDlcItems"`
		Unsearchable        bool          `json:"unsearchable"`
	} `json:"items"`
	CurrencyCode          string `json:"currencyCode"`
	CurrentPrice          int    `json:"currentPrice"`
	Price                 int    `json:"price"`
	BasePrice             int    `json:"basePrice"`
	BasePriceCurrencyCode string `json:"basePriceCurrencyCode"`
	RecurringPrice        int    `json:"recurringPrice"`
	FreeDays              int    `json:"freeDays"`
	MaxBillingCycles      int    `json:"maxBillingCycles"`
	Seller                struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"seller"`
	ViewableDate                       time.Time     `json:"viewableDate"`
	EffectiveDate                      time.Time     `json:"effectiveDate"`
	VatIncluded                        bool          `json:"vatIncluded"`
	IsCodeRedemptionOnly               bool          `json:"isCodeRedemptionOnly"`
	IsFeatured                         bool          `json:"isFeatured"`
	TaxSkuID                           string        `json:"taxSkuId"`
	MerchantGroup                      string        `json:"merchantGroup"`
	PriceTier                          string        `json:"priceTier"`
	URLSlug                            string        `json:"urlSlug"`
	RoleNamesToGrant                   []interface{} `json:"roleNamesToGrant"`
	Tags                               []interface{} `json:"tags"`
	PurchaseLimit                      int           `json:"purchaseLimit"`
	IgnoreOrder                        bool          `json:"ignoreOrder"`
	FulfillToGroup                     bool          `json:"fulfillToGroup"`
	FraudItemType                      string        `json:"fraudItemType"`
	ShareRevenue                       bool          `json:"shareRevenue"`
	OfferType                          string        `json:"offerType"`
	Unsearchable                       bool          `json:"unsearchable"`
	ReleaseOffer                       string        `json:"releaseOffer"`
	Title4Sort                         string        `json:"title4Sort"`
	SelfRefundable                     bool          `json:"selfRefundable"`
	RefundType                         string        `json:"refundType"`
	VisibilityType                     string        `json:"visibilityType"`
	CurrencyDecimals                   int           `json:"currencyDecimals"`
	AllowPurchaseForPartialOwned       bool          `json:"allowPurchaseForPartialOwned"`
	ShareRevenueWithUnderageAffiliates bool          `json:"shareRevenueWithUnderageAffiliates"`
	PlatformWhitelist                  []interface{} `json:"platformWhitelist"`
	PlatformBlacklist                  []interface{} `json:"platformBlacklist"`
	PartialItemPrerequisiteCheck       bool          `json:"partialItemPrerequisiteCheck"`
}

/*
{
  "app": "fortnite",
  "serverDate": "2023-04-15T00:49:58.443Z",
  "overridePropertiesVersion": "unknown",
  "cln": "25029190",
  "build": "255",
  "moduleName": "Fortnite-Core",
  "buildDate": "2023-04-13T18:00:01.364Z",
  "version": "24.20",
  "branch": "Release-24.20",
  "modules": {
    "Epic-LightSwitch-AccessControlCore": {
      "cln": "24565549",
      "build": "b2144",
      "buildDate": "2023-03-08T20:12:52.378Z",
      "version": "1.0.0",
      "branch": "trunk"
    },
    "epic-xmpp-api-v1-base": {
      "cln": "5131a23c1470acbd9c94fae695ef7d899c1a41d6",
      "build": "b3595",
      "buildDate": "2019-07-30T09:11:06.587Z",
      "version": "0.0.1",
      "branch": "master"
    },
    "epic-common-core": {
      "cln": "b3121b2c45b97e4ac06c789b052dc39f00a992fb",
      "build": "b1202",
      "buildDate": "2023-03-13T17:20:46.525Z",
      "version": "5.0.0.20230313171954",
      "branch": "master"
    }
  }
}
*/

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
