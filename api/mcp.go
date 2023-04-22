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
	RVN             int64              `json:"rvn"`
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

type CampaignProfile struct {
	ProfileRevision            int    `json:"profileRevision"`
	ProfileID                  string `json:"profileId"`
	ProfileChangesBaseRevision int    `json:"profileChangesBaseRevision"`
	ProfileChanges             []struct {
		ChangeType string              `json:"changeType"`
		Profile    CampaignProfileData `json:"profile"`
	} `json:"profileChanges"`
	ProfileCommandRevision int    `json:"profileCommandRevision"`
	ServerTime             string `json:"serverTime"`
	ResponseVersion        int    `json:"responseVersion"`
	Notifications          []struct {
		Type         string `json:"type"`
		Primary      bool   `json:"primary"`
		DaysLoggedIn int    `json:"daysLoggedIn"`
		Items        []struct {
			ItemType    string `json:"itemType"`
			ItemGuid    string `json:"itemGuid"`
			ItemProfile string `json:"itemProfile"`
			Quantity    int    `json:"quantity"`
		} `json:"items"`
	} `json:"notifications"`
}

type CampaignProfileData struct {
	Created         string                               `json:"created"`
	Updated         string                               `json:"updated"`
	RVN             int                                  `json:"rvn"`
	WipeNumber      int                                  `json:"wipeNumber"`
	AccountID       string                               `json:"accountId"`
	ProfileID       string                               `json:"profileId"`
	Version         string                               `json:"version"`
	CommandRevision int                                  `json:"commandRevision"`
	Items           map[string]*CampaignProfileItemEntry `json:"items"`
	Stats           CampaignProfileStats                 `json:"stats"`
}

type CampaignProfileItemEntry struct {
	TemplateID string `json:"templateId"`
	Quantity   int    `json:"quantity"`
	Attributes struct {
		QuestState string `json:"quest_state"`
	} `json:"attributes"`
}

type CampaignProfileStats struct {
	Attributes struct {
		NodeCosts struct {
			HomebaseNodeDefaultPage struct {
				TokenHomebasepoints int `json:"Token:homebasepoints"`
			} `json:"homebase_node_default_page"`
			ResearchNodeDefaultPage map[string]int `json:"research_node_default_page"`
		} `json:"node_costs"`
		MissionAlertRedemptionRecord struct {
			ClaimData []struct {
				MissionAlertID         string `json:"missionAlertId"`
				RedemptionDateUtc      string `json:"redemptionDateUtc"`
				EvictClaimDataAfterUtc string `json:"evictClaimDataAfterUtc"`
			} `json:"claimData"`
		} `json:"mission_alert_redemption_record"`
		ClientSettings struct {
			PinnedQuestInstances []interface{} `json:"pinnedQuestInstances"`
		} `json:"client_settings"`
		ResearchLevels struct {
			Fortitude  int `json:"fortitude"`
			Offense    int `json:"offense"`
			Resistance int `json:"resistance"`
			Technology int `json:"technology"`
		} `json:"research_levels"`
		Level               int           `json:"level"`
		SelectedHeroLoadout string        `json:"selected_hero_loadout"`
		Loadouts            []interface{} `json:"loadouts"`
		CollectionBook      struct {
			MaxBookXpLevelAchieved int `json:"maxBookXpLevelAchieved"`
		} `json:"collection_book"`
		LatentXPMarker   string `json:"latent_xp_marker"`
		MfaRewardClaimed bool   `json:"mfa_reward_claimed"`
		QuestManager     struct {
			DailyLoginInterval string `json:"dailyLoginInterval"`
			DailyQuestRerolls  int    `json:"dailyQuestRerolls"`
			QuestPoolStats     struct {
				PoolStats []struct {
					PoolName         string   `json:"poolName"`
					NextRefresh      string   `json:"nextRefresh"`
					RerollsRemaining int      `json:"rerollsRemaining"`
					QuestHistory     []string `json:"questHistory"`
				} `json:"poolStats"`
				DailyLoginInterval string `json:"dailyLoginInterval"`
				PoolLockouts       struct {
					PoolLockouts []struct {
						LockoutName string `json:"lockoutName"`
					} `json:"poolLockouts"`
				} `json:"poolLockouts"`
			} `json:"questPoolStats"`
		} `json:"quest_manager"`
		LegacyResearchPointsSpent int `json:"legacy_research_points_spent"`
		GameplayStats             []struct {
			StatName  string `json:"statName"`
			StatValue int    `json:"statValue"`
		} `json:"gameplay_stats"`
		EventCurrency struct {
			TemplateID string  `json:"templateId"`
			CF         float64 `json:"cf"`
		} `json:"event_currency"`
		MatchesPlayed int           `json:"matches_played"`
		ModeLoadouts  []interface{} `json:"mode_loadouts"`
		DailyRewards  struct {
			NextDefaultReward   int    `json:"nextDefaultReward"`
			TotalDaysLoggedIn   int    `json:"totalDaysLoggedIn"`
			LastClaimDate       string `json:"lastClaimDate"`
			AdditionalSchedules struct {
				Founderspackdailyrewardtoken struct {
					RewardsClaimed int  `json:"rewardsClaimed"`
					ClaimedToday   bool `json:"claimedToday"`
				} `json:"founderspackdailyrewardtoken"`
			} `json:"additionalSchedules"`
		} `json:"daily_rewards"`
		LastAppliedLoadout string `json:"last_applied_loadout"`
		XP                 int    `json:"xp"`
		PacksGranted       int    `json:"packs_granted"`
	} `json:"attributes"`
}

type CommonCoreProfile struct {
	ProfileRevision            int    `json:"profileRevision"`
	ProfileId                  string `json:"profileId"`
	ProfileChangesBaseRevision int    `json:"profileChangesBaseRevision"`
	ProfileCommandRevision     int    `json:"profileCommandRevision"`
	ServerTime                 string `json:"serverTime"`
	ResponseVersion            int    `json:"responseVersion"`
	ProfileChanges             []struct {
		ChangeType string                `json:"changeType"`
		Profile    CommonCoreProfileData `json:"profile"`
	} `json:"profileChanges"`
}

type CommonCoreProfileData struct {
	Created         string                                 `json:"created"`
	Updated         string                                 `json:"updated"`
	RVN             int                                    `json:"rvn"`
	WipeNumber      int                                    `json:"wipeNumber"`
	AccountID       string                                 `json:"accountId"`
	ProfileID       string                                 `json:"profileId"`
	Version         string                                 `json:"version"`
	CommandRevision int                                    `json:"commandRevision"`
	Items           map[string]*CommonCoreProfileItemEntry `json:"items"`
	Stats           CommonCoreProfileStats                 `json:"stats"`
}

type CommonCoreProfileItemEntry struct {
	TemplateID string `json:"templateId"`
	Quantity   int    `json:"quantity"`
	Attributes struct {
		Platform string `json:"platform"`
	} `json:"attributes"`
}

type CommonCoreProfileStats struct {
	Attributes struct {
		SurveyData struct {
			AllSurveysMetadata struct {
				NumTimesCompleted int    `json:"numTimesCompleted"`
				LastTimeCompleted string `json:"lastTimeCompleted"`
			} `json:"allSurveysMetadata"`
			Metadata map[string]struct {
				NumTimesCompleted int    `json:"numTimesCompleted"`
				LastTimeCompleted string `json:"lastTimeCompleted"`
			} `json:"metadata"`
		} `json:"survey_data"`
		IntroGamePlayed    bool `json:"intro_game_played"`
		MtxPurchaseHistory struct {
			RefundsUsed               int    `json:"refundsUsed"`
			RefundCredits             int    `json:"refundCredits"`
			TokenRefreshReferenceTime string `json:"tokenRefreshReferenceTime"`
			Purchases                 []struct {
				PurchaseID         string        `json:"purchaseId"`
				OfferID            string        `json:"offerId"`
				PurchaseDate       string        `json:"purchaseDate"`
				FreeRefundEligible bool          `json:"freeRefundEligible"`
				Fulfillments       []interface{} `json:"fulfillments"`
				LootResult         []struct {
					ItemType    string `json:"itemType"`
					ItemGUID    string `json:"itemGuid"`
					ItemProfile string `json:"itemProfile"`
					Quantity    int    `json:"quantity"`
				} `json:"lootResult"`
				TotalMtxPaid int `json:"totalMtxPaid"`
				Metadata     struct {
				} `json:"metadata"`
				GameContext string `json:"gameContext"`
				RefundDate  string `json:"refundDate,omitempty"`
				UndoTimeout string `json:"undoTimeout,omitempty"`
			} `json:"purchases"`
		} `json:"mtx_purchase_history"`
		UndoCooldowns []struct {
			OfferID         string `json:"offerId"`
			CooldownExpires string `json:"cooldownExpires"`
		} `json:"undo_cooldowns"`
		MtxAffiliateSetTime string `json:"mtx_affiliate_set_time"`
		CurrentMtxPlatform  string `json:"current_mtx_platform"`
		MtxAffiliate        string `json:"mtx_affiliate"`
		WeeklyPurchases     struct {
			LastInterval string         `json:"lastInterval"`
			PurchaseList map[string]int `json:"purchaseList"`
		} `json:"weekly_purchases"`
		DailyPurchases struct {
			LastInterval string         `json:"lastInterval"`
			PurchaseList map[string]int `json:"purchaseList"`
		} `json:"daily_purchases"`
		InAppPurchases struct {
			Receipts          []string       `json:"receipts"`
			IgnoredReceipts   []interface{}  `json:"ignoredReceipts"`
			FulfillmentCounts map[string]int `json:"fulfillmentCounts"`
			RefreshTimers     map[string]struct {
				NextEntitlementRefresh string `json:"nextEntitlementRefresh"`
			} `json:"refreshTimers"`
		} `json:"in_app_purchases"`
		ForcedIntroPlayed  string `json:"forced_intro_played"`
		RmtPurchaseHistory struct {
			Purchases []struct {
				FulfillmentID string `json:"fulfillmentId"`
				PurchaseDate  string `json:"purchaseDate"`
				LootResult    []struct {
					ItemType    string `json:"itemType"`
					ItemGUID    string `json:"itemGuid"`
					ItemProfile string `json:"itemProfile"`
					Attributes  struct {
						Platform string `json:"platform"`
					} `json:"attributes"`
					Quantity int `json:"quantity"`
				} `json:"lootResult"`
			} `json:"purchases"`
		} `json:"rmt_purchase_history"`
		UndoTimeout      string `json:"undo_timeout"`
		MonthlyPurchases struct {
			LastInterval string         `json:"lastInterval"`
			PurchaseList map[string]int `json:"purchaseList"`
		} `json:"monthly_purchases"`
		AllowedToSendGifts    bool   `json:"allowed_to_send_gifts"`
		MfaEnabled            bool   `json:"mfa_enabled"`
		AllowedToReceiveGifts bool   `json:"allowed_to_receive_gifts"`
		MtxAffiliateID        string `json:"mtx_affiliate_id"`
		GiftHistory           struct {
			NumSent      int               `json:"num_sent"`
			SentTo       map[string]string `json:"sentTo"`
			NumReceived  int               `json:"num_received"`
			ReceivedFrom map[string]string `json:"receivedFrom"`
			Gifts        []struct {
				Date        string `json:"date"`
				OfferID     string `json:"offerId"`
				ToAccountID string `json:"toAccountId"`
			} `json:"gifts"`
		} `json:"gift_history"`
	} `json:"attributes"`
}
