package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type PartyMetaUpdateAthenaCosmeticStat struct {
	StatName  string `json:"statName"`
	StatValue int    `json:"statValue"`
}

// type PartyMetaUpdateAthenaCosmeticLoadout struct {
// 	CharacterDef  string                              `json:"characterDef"`
// 	CharacterEKey string                              `json:"characterEKey"`
// 	BackpackDef   string                              `json:"backpackDef"`
// 	BackpackEKey  string                              `json:"backpackEKey"`
// 	PickaxeDef    string                              `json:"pickaxeDef"`
// 	PickaxeEKey   string                              `json:"pickaxeEKey"`
// 	ContrailDef   string                              `json:"contrailDef"`
// 	ContrailEKey  string                              `json:"contrailEKey"`
// 	Scratchpad    []string                            `json:"scratchpad"`
// 	CosmeticStats []PartyMetaUpdateAthenaCosmeticStat `json:"cosmeticStats"`
// }

type PartyMetaUpdate struct {
	Delete   []string            `json:"delete"`
	Revision int                 `json:"revision"`
	Update   PartyMetaUpdateData `json:"update"`
}

type PartyMetaUpdateData struct {
	DefaultAthenaCosmeticLoadoutJ PartyMetaUpdateDefaultAthenaCosmeticLoadoutJ `json:"Default:AthenaCosmeticLoadout_j"`
}

type PartyMetaUpdateDefaultAthenaCosmeticLoadout struct {
	CharacterDef  string   `json:"characterDef"`
	CharacterEKey string   `json:"characterEKey"`
	BackpackDef   string   `json:"backpackDef"`
	BackpackEKey  string   `json:"backpackEKey"`
	PickaxeDef    string   `json:"pickaxeDef"`
	PickaxeEKey   string   `json:"pickaxeEKey"`
	ContrailDef   string   `json:"contrailDef"`
	ContrailEKey  string   `json:"contrailEKey"`
	Scratchpad    []string `json:"scratchpad"`
	CosmeticStats []PartyMetaUpdateAthenaCosmeticStat
}

type PartyMetaUpdateDefaultAthenaCosmeticLoadoutJ struct {
	AthenaCosmeticLoadout PartyMetaUpdateDefaultAthenaCosmeticLoadout `json:"AthenaCosmeticLoadout"`
}

// type PartyMetaUpdate struct {
// 	DefaultAthenaCosmeticLoadoutJ PartyMetaUpdateDefaultAthenaCosmeticLoadoutJ `json:"Default:AthenaCosmeticLoadout_j"`
// }

/* PATCH https://party-service-prod.ol.epicgames.com/party/api/v1/Fortnite/parties/{PartyID}/members/{Account_ID}/meta
{
	"delete": [],
	"revision": 4,
	"update": {
		"Default:AthenaCosmeticLoadout_j": "{
			'AthenaCosmeticLoadout': {
				'characterDef': '/Game/Athena/Items/Cosmetics/Characters/CID_342_Athena_Commando_M_StreetRacerMetallic.CID_342_Athena_Commando_M_StreetRacerMetallic',
				'characterEKey': '',
				'backpackDef': '/Game/Athena/Items/Cosmetics/Backpacks/BID_610_ElasticHologram.BID_610_ElasticHologram',
				'backpackEKey': '',
				'pickaxeDef': '/Game/Athena/Items/Cosmetics/Pickaxes/Pickaxe_ID_035_Prismatic.Pickaxe_ID_035_Prismatic',
				'pickaxeEKey': '',
				'contrailDef': '/Game/Athena/Items/Cosmetics/Contrails/Trails_ID_019_PSBurnout.Trails_ID_019_PSBurnout',
				'contrailEKey': '',
				'scratchpad': [],
				'cosmeticStats': [
					{
						'statName': 'TotalVictoryCrowns',
						'statValue': 999
					},
					{
						'statName': 'TotalRoyalRoyales',
						'statValue': 999
					},
					{
						'statName':
						'HasCrown',
						'statValue': 1
					}
				]
			}
		}"
	}
}
*/

/*
map[
	current:[
		map[
			applicants:[]
			config:map[
				discoverability:INVITED_ONLY
				intention_ttl:60
				invite_ttl:14400
				join_confirmation:true
				joinability:INVITE_AND_FORMER
				max_size:16
				sub_type:default
				type:DEFAULT
			]
			created_at:2023-04-15T11:29:16.052Z
			id:be57c484df484a38acca7f5ccd6d447d
			intentions:[]
			invites:[]
			members:[
				map[
					account_id:3de654385bf64a9686614cdc999d2d21
					connections:[
						map[
							account_pl:xbl
							account_pl_dn:Distrust7656
							connected_at:2023-04-15T11:29:15.533Z
							id:3de654385bf64a9686614cdc999d2d21@prod.ol.epicgames.com/V2:Fortnite:XSX:2535416849881588:F567ABF043B021ADBEBF0AA03992471D
							meta:map[urn:epic:conn:platform:dn_s:Distrust7656 urn:epic:conn:platform_s:XSX]
							updated_at:2023-04-15T11:29:16.052Z
							yield_leadership:false
						]
					]
					joined_at:2023-04-15T11:29:16.057Z
					meta:map[
						Default:ArbitraryCustomDataStore_j:{
							"ArbitraryCustomDataStore":[]
						}
						Default:AthenaBannerInfo_j:{
							"AthenaBannerInfo":{
								"bannerIconId":"brseason01",
								"bannerColorId":"defaultcolor39",
								"seasonLevel":1
							}
						}
						Default:AthenaCosmeticLoadoutVariants_j:{
							"AthenaCosmeticLoadoutVariants":{
								"vL":{
									"athenaCharacter":{
										"i":[
											{
												"c":"Material",
												"v":"Mat3",
												"dE":0
											},
											{
											"c":"Parts",
											"v":"Stage1",
											"dE":0
											},
											{
												"c":"Progressive",
												"v":"Stage1",
												"dE":0
											}
										]
									}
								},
								"fT":false
							}
						}
						Default:AthenaCosmeticLoadout_j:{
							"AthenaCosmeticLoadout":{
								"characterDef":"/Game/Athena/Items/Cosmetics/Characters/CID_A_371_Athena_Commando_F_Cadet.CID_A_371_Athena_Commando_F_Cadet",
								"characterEKey":"",
								"backpackDef":"None",
								"backpackEKey":"",
								"pickaxeDef":"/Game/Athena/Items/Cosmetics/Pickaxes/Pickaxe_ID_749_GimmickMale_5C033.Pickaxe_ID_749_GimmickMale_5C033",
								"pickaxeEKey":"",
								"contrailDef":"/Game/Athena/Items/Cosmetics/Contrails/Trails_ID_001_Disco.Trails_ID_001_Disco",
								"contrailEKey":"",
								"scratchpad":[],
								"cosmeticStats":[
									{"statName":"TotalVictoryCrowns","statValue":0},
									{"statName":"TotalRoyalRoyales","statValue":0},
									{"statName":"HasCrown","statValue":0}
								]
							}
						}
						Default:BattlePassInfo_j:{"BattlePassInfo":{"bHasPurchasedPass":false,"passLevel":1,"selfBoostXp":0,"friendBoostXp":0}} Default:CampaignHero_j:{"CampaignHero":{"heroItemInstanceId":"","heroType":"/Game/Athena/Heroes/HID_A_371_Athena_Commando_F_Cadet.HID_A_371_Athena_Commando_F_Cadet"}}
						Default:CampaignInfo_j:{"CampaignInfo":{"matchmakingLevel":0,"zoneInstanceId":"","homeBaseVersion":1}}
						Default:CrossplayPreference_s:OptedIn Default:DownloadOnDemandProgress_d:0.000000
						Default:FeatDefinition_s:None
						Default:FortCommonMatchmakingData_j:{"FortCommonMatchmakingData":{"request":{"linkId":{"mnemonic":"","version":-1},"requester":"INVALID","version":0},"response":"NONE"}}
						Default:FrontEndMapMarker_j:{"FrontEndMapMarker":{"markerLocation":{"x":0,"y":0},"bIsSet":false}} Default:FrontendEmote_j:{"FrontendEmote":{"emoteItemDef":"None","emoteEKey":"","emoteSection":-1}}
						Default:JoinInProgressData_j:{"JoinInProgressData":{"request":{"target":"INVALID","time":0},"responses":[]}}
						Default:JoinMethod_s:Creation Default:LobbyState_j:{"LobbyState":{"inGameReadyCheckStatus":"None","gameReadiness":"NotReady","readyInputType":"Count","currentInputType":"Gamepad","hiddenMatchmakingDelayMax":0,"hasPreloadedAthena":false}}
						Default:MemberSquadAssignmentRequest_j:{"MemberSquadAssignmentRequest":{"startingAbsoluteIdx":-1,"targetAbsoluteIdx":-1,"swapTargetMemberId":"INVALID","version":0}}
						Default:NumAthenaPlayersLeft_U:0 Default:PackedState_j:{"PackedState":{"subGame":"Athena","location":"PreLobby","gameMode":"None","voiceChatStatus":"Enabled","hasCompletedSTWTutorial":true,"hasPurchasedSTW":true,"platformSupportsSTW":true,"bDownloadOnDemandActive":false,"bIsPartyLFG":false,"bShouldRecordPartyChannel":false}}
						Default:PlatformData_j:{"PlatformData":{"platform":{"platformDescription":{"name":"HELIOS","platformType":"CONSOLE","onlineSubsystem":"GDK","sessionType":"MPSD","externalAccountType":"xbl","crossplayPool":"xbl"}},"uniqueId":"GDK:2535416849881588","sessionId":""}}
						Default:SharedQuests_j:{"SharedQuests":{"bcktMap":{"BR":{"qsts":[{"qst":"Quest:quest_s24_dailyquest_damageopponents","objs":[0]},{"qst":"Quest:quest_s24_dailyquest_distanceskating","objs":[0]},{"qst":"Quest:quest_s24_dailyquest_healvegmeat","objs":[0]}]},"BRBackup":{"qsts":[{"qst":"Quest:quest_s24_dailyquestbackup_02","objs":[0]},{"qst":"Quest:quest_s24_dailyquestbackup_01","objs":[0]}]}},"pndQst":"Quest:quest_s24_boosterpack_q01"}}
						Default:SpectateInfo_j:{"SpectateInfo":{"gameSessionId":"","gameSessionKey":""}}
						Default:UtcTimeStartedMatchAthena_s:0001-01-01T00:00:00.000Z
						Default:bIsPartyUsingPartySignal_b:false urn:epic:member:dn_s:Distrust7656] revision:2 role:CAPTAIN updated_at:2023-04-15T11:34:16.919Z]] meta:map[
							Default:ActivityName_s:
							Default:ActivityType_s:Undefined
							Default:AllowJoinInProgress_b:false
							Default:AthenaPrivateMatch_b:false Default:AthenaSquadFill_b:true
							Default:CampaignInfo_j:{"CampaignInfo":{"lobbyConnectionStarted":false,"matchmakingResult":"NoResults","matchmakingState":"NotMatchmaking","sessionIsCriticalMission":false,"zoneTileIndex":-1,"theaterId":"","tileStates":{"tileStates":[],"numSetBits":0}}}
							Default:CreativeDiscoverySurfaceRevisions_j:{"CreativeDiscoverySurfaceRevisions":[]}
							Default:CreativePortalCountdownStartTime_s:0001-01-01T00:00:00.000Z
							Default:CurrentRegionId_s:NAE Default:CustomMatchKey_s:
							Default:FortCommonMatchmakingData_j:{"FortCommonMatchmakingData":{"current":{"linkId":{"mnemonic":"playlist_defaultsolo","version":-1},"requester":"INVALID","version":12194},"participantData":{"requested":{"linkId":{"mnemonic":"playlist_defaultsolo","version":-1},"requester":"INVALID","version":12194},"broadcast":"ReadyForRequests"}}}
							Default:GameSessionKey_s:
							Default:LFGTime_s:0001-01-01T00:00:00.000Z
							Default:PartyIsJoinedInProgress_b:false
							Default:PartyMatchmakingInfo_j:{"PartyMatchmakingInfo":{"buildId":-1,"hotfixVersion":-1,"regionId":"","playlistName":"None","playlistRevision":0,"tournamentId":"","eventWindowId":"","linkCode":""}}
							Default:PartyState_s:BattleRoyaleView Default:PlatformSessions_j:{"PlatformSessions":[{"sessionType":"MPSD","sessionId":"","ownerPrimaryId":"3de654385bf64a9686614cdc999d2d21"}]}
							Default:PlaylistData_j:{"PlaylistData":{"playlistName":"Playlist_DefaultSolo","tournamentId":"","eventWindowId":"","regionId":"NONE","linkId":{"mnemonic":"playlist_defaultsolo","version":-1},"bGracefullyUpgraded":false,"matchmakingRulePreset":"RespectParties"}}
							Default:PrimaryGameSessionId_s: Default:PrivacySettings_j:{"PrivacySettings":{"partyType":"Private","partyInviteRestriction":"AnyMember","bOnlyLeaderFriendsCanJoin":false}}
							Default:RawSquadAssignments_j:{"RawSquadAssignments":[{"memberId":"3de654385bf64a9686614cdc999d2d21","absoluteMemberIdx":0}]}
							Default:ZoneInstanceId_s: urn:epic:cfg:accepting-members_b:false urn:epic:cfg:build-id_s:1:3:24849278 urn:epic:cfg:can-join_b:true urn:epic:cfg:chat-enabled_b:true urn:epic:cfg:invite-perm_s:Anyone urn:epic:cfg:join-request-action_s:Manual urn:epic:cfg:not-accepting-members-reason_i:7 urn:epic:cfg:party-type-id_s:default urn:epic:cfg:presence-perm_s:Noone] revision:2 updated_at:2023-04-15T11:29:24.702Z]] invites:[] pending:[] pings:[]
							]
*/

type Party struct {
	Current []struct {
		ID        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Config    struct {
			Type             string `json:"type"`
			Joinability      string `json:"joinability"`
			Discoverability  string `json:"discoverability"`
			SubType          string `json:"sub_type"`
			MaxSize          int    `json:"max_size"`
			InviteTTL        int    `json:"invite_ttl"`
			JoinConfirmation bool   `json:"join_confirmation"`
			IntentionTTL     int    `json:"intention_ttl"`
		} `json:"config"`
		Members []struct {
			AccountID string `json:"account_id"`
			Meta      struct {
				DefaultFrontEndMapMarkerJ             string    `json:"Default:FrontEndMapMarker_j"`
				DefaultArbitraryCustomDataStoreJ      string    `json:"Default:ArbitraryCustomDataStore_j"`
				DefaultNumAthenaPlayersLeftU          string    `json:"Default:NumAthenaPlayersLeft_U"`
				DefaultAthenaCosmeticLoadoutJ         string    `json:"Default:AthenaCosmeticLoadout_j"`
				UrnEpicMemberDnS                      string    `json:"urn:epic:member:dn_s"`
				DefaultBattlePassInfoJ                string    `json:"Default:BattlePassInfo_j"`
				DefaultCampaignInfoJ                  string    `json:"Default:CampaignInfo_j"`
				DefaultPackedStateJ                   string    `json:"Default:PackedState_j"`
				DefaultFeatDefinitionS                string    `json:"Default:FeatDefinition_s"`
				DefaultCampaignHeroJ                  string    `json:"Default:CampaignHero_j"`
				DefaultFortCommonMatchmakingDataJ     string    `json:"Default:FortCommonMatchmakingData_j"`
				DefaultSpectateInfoJ                  string    `json:"Default:SpectateInfo_j"`
				DefaultFrontendEmoteJ                 string    `json:"Default:FrontendEmote_j"`
				DefaultBIsPartyUsingPartySignalB      string    `json:"Default:bIsPartyUsingPartySignal_b"`
				DefaultAthenaCosmeticLoadoutVariantsJ string    `json:"Default:AthenaCosmeticLoadoutVariants_j"`
				DefaultMemberSquadAssignmentRequestJ  string    `json:"Default:MemberSquadAssignmentRequest_j"`
				DefaultCrossplayPreferenceS           string    `json:"Default:CrossplayPreference_s"`
				DefaultAthenaBannerInfoJ              string    `json:"Default:AthenaBannerInfo_j"`
				DefaultPlatformDataJ                  string    `json:"Default:PlatformData_j"`
				DefaultDownloadOnDemandProgressD      string    `json:"Default:DownloadOnDemandProgress_d"`
				DefaultJoinMethodS                    string    `json:"Default:JoinMethod_s"`
				DefaultSharedQuestsJ                  string    `json:"Default:SharedQuests_j"`
				DefaultUtcTimeStartedMatchAthenaS     time.Time `json:"Default:UtcTimeStartedMatchAthena_s"`
				DefaultLobbyStateJ                    string    `json:"Default:LobbyState_j"`
				DefaultJoinInProgressDataJ            string    `json:"Default:JoinInProgressData_j"`
			} `json:"meta"`
			Connections []struct {
				ID              string    `json:"id"`
				ConnectedAt     time.Time `json:"connected_at"`
				UpdatedAt       time.Time `json:"updated_at"`
				AccountPl       string    `json:"account_pl"`
				AccountPlDn     string    `json:"account_pl_dn"`
				YieldLeadership bool      `json:"yield_leadership"`
				Meta            struct {
					UrnEpicConnPlatformS   string `json:"urn:epic:conn:platform_s"`
					UrnEpicConnPlatformDnS string `json:"urn:epic:conn:platform:dn_s"`
				} `json:"meta"`
			} `json:"connections"`
			Revision  int       `json:"revision"`
			UpdatedAt time.Time `json:"updated_at"`
			JoinedAt  time.Time `json:"joined_at"`
			Role      string    `json:"role"`
		} `json:"members"`
		Applicants []interface{} `json:"applicants"`
		Meta       struct {
			DefaultAllowJoinInProgressB               string    `json:"Default:AllowJoinInProgress_b"`
			DefaultPartyIsJoinedInProgressB           string    `json:"Default:PartyIsJoinedInProgress_b"`
			DefaultRawSquadAssignmentsJ               string    `json:"Default:RawSquadAssignments_j"`
			DefaultPlaylistDataJ                      string    `json:"Default:PlaylistData_j"`
			DefaultPartyStateS                        string    `json:"Default:PartyState_s"`
			DefaultCampaignInfoJ                      string    `json:"Default:CampaignInfo_j"`
			UrnEpicCfgPartyTypeIDS                    string    `json:"urn:epic:cfg:party-type-id_s"`
			DefaultFortCommonMatchmakingDataJ         string    `json:"Default:FortCommonMatchmakingData_j"`
			UrnEpicCfgBuildIDS                        string    `json:"urn:epic:cfg:build-id_s"`
			UrnEpicCfgPresencePermS                   string    `json:"urn:epic:cfg:presence-perm_s"`
			DefaultCustomMatchKeyS                    string    `json:"Default:CustomMatchKey_s"`
			UrnEpicCfgAcceptingMembersB               string    `json:"urn:epic:cfg:accepting-members_b"`
			DefaultLFGTimeS                           time.Time `json:"Default:LFGTime_s"`
			UrnEpicCfgJoinRequestActionS              string    `json:"urn:epic:cfg:join-request-action_s"`
			DefaultPrimaryGameSessionIDS              string    `json:"Default:PrimaryGameSessionId_s"`
			DefaultAthenaPrivateMatchB                string    `json:"Default:AthenaPrivateMatch_b"`
			UrnEpicCfgInvitePermS                     string    `json:"urn:epic:cfg:invite-perm_s"`
			DefaultPrivacySettingsJ                   string    `json:"Default:PrivacySettings_j"`
			DefaultActivityTypeS                      string    `json:"Default:ActivityType_s"`
			UrnEpicCfgNotAcceptingMembersReasonI      string    `json:"urn:epic:cfg:not-accepting-members-reason_i"`
			DefaultCreativeDiscoverySurfaceRevisionsJ string    `json:"Default:CreativeDiscoverySurfaceRevisions_j"`
			DefaultGameSessionKeyS                    string    `json:"Default:GameSessionKey_s"`
			UrnEpicCfgChatEnabledB                    string    `json:"urn:epic:cfg:chat-enabled_b"`
			DefaultAthenaSquadFillB                   string    `json:"Default:AthenaSquadFill_b"`
			DefaultZoneInstanceIDS                    string    `json:"Default:ZoneInstanceId_s"`
			DefaultPartyMatchmakingInfoJ              string    `json:"Default:PartyMatchmakingInfo_j"`
			DefaultPlatformSessionsJ                  string    `json:"Default:PlatformSessions_j"`
			UrnEpicCfgCanJoinB                        string    `json:"urn:epic:cfg:can-join_b"`
			DefaultActivityNameS                      string    `json:"Default:ActivityName_s"`
			DefaultCreativePortalCountdownStartTimeS  time.Time `json:"Default:CreativePortalCountdownStartTime_s"`
			DefaultCurrentRegionIDS                   string    `json:"Default:CurrentRegionId_s"`
		} `json:"meta"`
		Invites    []interface{} `json:"invites"`
		Revision   int           `json:"revision"`
		Intentions []interface{} `json:"intentions"`
	} `json:"current"`
	Pending []interface{} `json:"pending"`
	Invites []interface{} `json:"invites"`
	Pings   []interface{} `json:"pings"`
}

func (c EpicClient) FetchParty(credentials UserCredentialsResponse) (Party, error) {
	// https://party-service-prod.ol.epicgames.com/party/api/v1/:namespace/user/:accountId
	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+credentials.AccessToken)

	resp, err := c.Request("GET", fmt.Sprintf("https://party-service-prod.ol.epicgames.com/party/api/v1/Fortnite/user/%s", credentials.AccountID), headers, "")
	if err != nil {
		return Party{}, err
	}

	defer resp.Body.Close()

	// log.Println(resp.Body)

	var res Party
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return Party{}, err
	}

	return res, nil
	// data, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return nil, err
	// }

	// return string(data), nil
}

func (c EpicClient) PartyMetaUpdate(credentials UserCredentialsResponse, partyID string, body string) (any, error) {
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Authorization", "Bearer "+credentials.AccessToken)

	resp, err := c.Request("PATCH", fmt.Sprintf("https://party-service-prod.ol.epicgames.com/party/api/v1/Fortnite/parties/%s/members/%s/meta", partyID, credentials.AccountID), headers, body)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var res any
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return MnemonicInfoResponse{}, err
	}

	return res, nil
}
