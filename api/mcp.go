package api

import (
	"blast/api/consts"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c EpicClient) ProfileOperationStr(credentials UserCredentialsResponse, operation string, profileId string, body string) (io.ReadCloser, error) {
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Authorization", "Bearer "+credentials.AccessToken)

	// https://fngw-mcp-gc-livefn.ol.epicgames.com/fortnite/api/game/v2/profile/:accountId/:route/:operation

	resp, err := c.Request("POST", fmt.Sprintf("%s/profile/%s/client/%s?profileId=%s&rvn=-1", consts.FORTNITE_GAME_BASE, credentials.AccountID, operation, profileId), headers, body)
	if err != nil {
		return nil, err
	}

	// defer resp.Body.Close()

	// bodyBytes, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return "", err
	// }

	return resp.Body, nil
}

func (c EpicClient) ProfileOperation(credentials UserCredentialsResponse, operation string, profileId string, body any) (any, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return c.ProfileOperationStr(credentials, operation, profileId, string(bodyBytes))
}

type AthenaProfile[T interface{}] struct {
	ProfileRevision            int    `json:"profileRevision"`
	ProfileID                  string `json:"profileId"`
	ProfileChangesBaseRevision int    `json:"profileChangesBaseRevision"`
	ProfileChanges             []struct {
		ChangeType string               `json:"changeType"`
		Profile    AthenaProfileData[T] `json:"profile"`
	} `json:"profileChanges"`
	ProfileCommandRevision int    `json:"profileCommandRevision"`
	ServerTime             string `json:"serverTime"`
	ResponseVersion        int    `json:"responseVersion"`
}

// type AthenaProfileItem interface {
// 	AthenaProfileLockerItem | AthenaProfile
// }

type AthenaProfileData[T interface{}] struct {
	Created         string             `json:"created"`
	Updated         string             `json:"updated"`
	Rvn             int64              `json:"rvn"`
	WipeNumber      int64              `json:"wipeNumber"`
	AccountID       string             `json:"accountId"`
	ProfileID       string             `json:"profileId"`
	Version         string             `json:"version"`
	Items           map[string]T       `json:"items"`
	Stats           AthenaProfileStats `json:"stats"`
	CommandRevision int64              `json:"commandRevision"`
}

// type E = AthenaProfileData | AthenaProfileLockerItem

type AthenaProfileLockerItem struct {
	TemplateID string `json:"templateId"`
	Attributes struct {
		LockerSlotsData struct {
			Slots map[string]struct {
				Items          []string `json:"items"`
				ActiveVariants []struct {
					Variants []AthenaVariant `json:"variants"`
				} `json:"activeVariants"`
			} `json:"slots"`
		} `json:"locker_slots_data"`
		UseCount            int64  `json:"use_count"`
		BannerIconTemplate  string `json:"banner_icon_template"`
		BannerColorTemplate string `json:"banner_color_template"`
		LockerName          string `json:"locker_name"`
	} `json:"attributes"`
	Quantity int64 `json:"quantity"`
}

/*
   "f30c612e-41c1-4454-a9dd-47a04ab2aad2" : {
     "templateId" : "CosmeticLocker:cosmeticlocker_athena",
     "attributes" : {
       "locker_slots_data" : {
         "slots" : {
           "LoadingScreen" : {
             "items" : [ "AthenaLoadingScreen:lsid_random" ],
             "activeVariants" : [ null ]
           }
         }
       },
       "use_count" : 0,
       "banner_icon_template" : "brseason01",
       "banner_color_template" : "defaultcolor39",
       "locker_name" : "PRESET 24"
     },
     "quantity" : 1
   },
*/

type AthenaVariant struct {
	Channel string   `json:"channel"`
	Active  string   `json:"active"`
	Owned   []string `json:"owned"`
}

type AthenaProfileStats struct {
	Attributes struct {
		UseRandomLoadout bool `json:"use_random_loadout"`
		PastSeasons      []struct {
			SeasonNumber    int64 `json:"seasonNumber"`
			NumWins         int64 `json:"numWins"`
			NumHighBracket  int64 `json:"numHighBracket"`
			NumLowBracket   int64 `json:"numLowBracket"`
			SeasonXp        int64 `json:"seasonXp"`
			SeasonLevel     int64 `json:"seasonLevel"`
			BookXP          int64 `json:"bookXp"`
			BookLevel       int64 `json:"bookLevel"`
			PurchasedVIP    bool  `json:"purchasedVIP"`
			NumRoyalRoyales int64 `json:"numRoyalRoyales"`
		} `json:"past_seasons"`
		SeasonMatchBoost          int64    `json:"season_match_boost"`
		Loadouts                  []string `json:"loadouts"`
		RestedXPOverflow          int64    `json:"rested_xp_overflow"`
		MFARewardClaimed          bool     `json:"mfa_reward_claimed"`
		LastXPInteraction         string   `json:"last_xp_interaction"`
		RestedXPGoldenPathGranted int64    `json:"rested_xp_golden_path_granted"`
		QuestManager              struct {
			DailyLoginInterval string `json:"dailyLoginInterval"`
			DailyQuestRerolls  int64  `json:"dailyQuestRerolls"`
		} `json:"quest_manager"`
		BookLevel         int64 `json:"book_level"`
		SeasonNum         int64 `json:"season_num"`
		SeasonUpdate      int64 `json:"season_update"`
		CreativeDynamicXP struct {
			Timespan          float64 `json:"timespan"`
			BucketXP          int64   `json:"bucketXp"`
			BankXP            int64   `json:"bankXp"`
			BankXpMult        float64 `json:"bankXpMult"`
			BoosterBucketXP   int64   `json:"boosterBucketXp"`
			BoosterXpMult     float64 `json:"boosterXpMult"`
			DailyExcessXPMult float64 `json:"dailyExcessXpMult"`
			CurrentDayXP      int64   `json:"currentDayXp"`
			CurrentDay        int64   `json:"currentDay"`
		} `json:"creative_dynamic_xp"`
		Season struct {
			NumWins        int64 `json:"numWins"`
			NumHighBracket int64 `json:"numHighBracket"`
			NumLowBracket  int64 `json:"numLowBracket"`
		} `json:"season"`
		Battlestars int64 `json:"battlestars"`
		Vote_data   struct {
			ElectionId  string `json:"electionId"`
			VoteHistory map[string]struct {
				VoteCount   int64  `json:"voteCount"`
				FirstVoteAt string `json:"firstVoteAt"`
				LastVoteAt  string `json:"lastVoteAt"`
			} `json:"voteHistory"`
			VotesRemaining  int64  `json:"votesRemaining"`
			LastVoteGranted string `json:"lastVoteGranted"`
		} `json:"vote_data"`
		BattlestarsSeasonTotal        int64  `json:"battlestars_season_total"`
		LifetimeWins                  int64  `json:"lifetime_wins"`
		PartyAssistQuest              string `json:"party_assist_quest"`
		BookPurchased                 bool   `json:"book_purchased"`
		PurchasedBattlePassTierOffers []struct {
			Id    string `json:"id"`
			Count int64  `json:"count"`
		} `json:"purchased_battle_pass_tier_offers"`
		RestedXPExchange       float64 `json:"rested_xp_exchange"`
		Level                  int64   `json:"level"`
		RestedXPMult           float64 `json:"rested_xp_mult"`
		AccountLevel           int64   `json:"accountLevel"`
		PinnedQuest            string  `json:"pinned_quest"`
		LastAppliedLoadout     string  `json:"last_applied_loadout"`
		XP                     int64   `json:"xp"`
		SeasonFriendMatchBoost int64   `json:"season_friend_match_boost"`
		PurchasedBPOffers      []struct {
			OfferId           string `json:"offerId"`
			BIsFreePassReward bool   `json:"bIsFreePassReward"`
			PurchaseDate      string `json:"purchaseDate"`
			LootResult        []struct {
				ItemType    string `json:"itemType"`
				ItemGuid    string `json:"itemGuid"`
				ItemProfile string `json:"itemProfile"`
				Attributes  struct {
					Platform string `json:"platform"`
				} `json:"attributes"`
				Quantity int64 `json:"quantity"`
			} `json:"lootResult"`
			CurrencyType      string `json:"currencyType"`
			TotalCurrencyPaid int64  `json:"totalCurrencyPaid"`
		} `json:"purchased_bp_offers"`
		LastMatchEndDatetime            string `json:"last_match_end_datetime"`
		LastSTWAccoladeTransferDatetime string `json:"last_stw_accolade_transfer_datetime"`
		MtxPurchaseHistoryCopy          []struct {
			PurchaseId   string   `json:"purchaseId"`
			PurchaseDate string   `json:"purchaseDate"`
			TemplateIds  []string `json:"templateIds"`
		} `json:"mtx_purchase_history_copy"`
	} `json:"attributes"`
}

// type AthenaProfileVictoryCrownProfileItem

// TODO
// "VictoryCrown:defaultvictorycrown": {
// 	"templateId": "VictoryCrown:defaultvictorycrown",
// 	"attributes": {
// 		"victory_crown_account_data": {
// 			"has_victory_crown": false,
// 			"data_is_valid_for_mcp": true,
// 			"total_victory_crowns_bestowed_count": 8,
// 			"total_royal_royales_achieved_count": 4
// 		},
// 		"level": 1
// 	},
// 	"quantity": 1
// }

// struct methods cannot use generic types, so this is seperated
// func ProfileOperation[T interface{}](c EpicClient, credentials UserCredentialsResponse, operation string, profileId string, body T) (any, error) {
// 	bodyBytes, err := json.Marshal(body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return c.ProfileOperationBasic(credentials, operation, profileId, string(bodyBytes))
// }

// func (c EpicClient) AddFavoriteMnemonic(credentials UserCredentialsResponse, mnemonic string) error {
// 	headers := http.Header{}
// 	headers.Set("Content-Type", "application/json")
// 	headers.Set("Authorization", "Bearer "+credentials.AccessToken)

// 	resp, err := c.Request("POST", fmt.Sprintf("https://fn-service-discovery-live-public.ogs.live.on.epicgames.com/api/v1/links/favorites/%s/%s", credentials.AccountID, mnemonic), headers, "{}")
// 	if err != nil {
// 		return err
// 	}

// 	if resp.StatusCode != http.StatusNoContent {
// 		return fmt.Errorf("invalid mnemonic")
// 	}

// 	return nil
// }

// export async function profileOperationRequest(client: Client, profileOperation: string, profileId: ProfileId, payload?: object) {
//     const { body, statusCode } = await httpRequest(client, `https://fortnite-public-service-prod11.ol.epicgames.com/fortnite/api/game/v2/profile/${client.accountId}/client/${profileOperation}?profileId=${profileId}&rvn=-1`, {
//         method: "POST",
//         headers: {
//             "Content-Type": "application/json"
//         },
//         body: JSON.stringify(payload ?? {})
//     });

//     return { body, statusCode };
// };
