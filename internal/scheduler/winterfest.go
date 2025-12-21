package scheduler

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/8h9x/fortgo"
	"github.com/8h9x/fortgo/fortnite"
	"github.com/8h9x/fortgo/request"
)

type RawItem[T any] struct {
	TemplateID string `json:"templateId"`
	Attributes T      `json:"attributes"`
	Quantity   int    `json:"quantity"`
}

type RewardKey struct {
	StaticKeyTemplateID string `json:"static_key_template_id"`
	UnlockKeysUsed      int    `json:"unlock_keys_used"`
	KeysGrantedToday    int    `json:"keys_granted_today"`
}

type RewardGraphItemAttributes struct {
	UnlockEpoch                   time.Time   `json:"unlock_epoch"`
	PlayerRandomSeed              int         `json:"player_random_seed"`
	RewardGraphPurchasedTimestamp int64       `json:"reward_graph_purchased_timestamp"`
	RewardGraphPurchased          bool        `json:"reward_graph_purchased"`
	RewardKeys                    []RewardKey `json:"reward_keys"`
}

type RewardGraphItem RawItem[RewardGraphItemAttributes]

type TokenItemAttributes struct {
	Level int `json:"level"`
}

type WinterfestGiftService struct{
	RewardGraphTemplateID string
}

func (s *WinterfestGiftService) Start() {
    log.Println("Starting WinterfestGiftService gift opening loop...") // TODO: discord channel alert

	mockSessions := []fortgo.Client{}

	for _, session := range mockSessions {
		err := s.OpenGift(session, "ERG.Node.A.1")
		if err != nil {
			log.Println(err) // TODO: discord alert to user
		}

		// TODO: update last open ts
	}
}

func (s *WinterfestGiftService) OpenGift(client fortgo.Client, rewardNodeID string) error {
	res, err := client.ClientQuestLogin("athena", "")
	if err != nil {
		return fmt.Errorf("WinterfestGiftService: unable to claim rewardNode: %s\n%s", rewardNodeID, err)
	}

	profileRes, err := request.ParseResponse[fortnite.Profile[fortnite.AthenaProfileStats, []any]](res)
	if err != nil {
		return fmt.Errorf("unable to parse ClientQuestLogin response: %s", err)
	}

	items := profileRes.Data.ProfileChanges[0].Profile.Items

	var rewardGraphID string
	var rewardGraphItem RewardGraphItem

	for key, value := range items {
		var item RewardGraphItem
		if err = json.Unmarshal(value, &item); err != nil {
			// not a skin; (you should probably add an additional check to ensure that it isn't some other type of error occurring); TODO: abstract this to a helper function that properly error checks and returns an empty state of the type passed if the type of data doesnt match
			continue
		}

		if strings.HasPrefix(item.TemplateID, s.RewardGraphTemplateID) {
			rewardGraphID = key
			rewardGraphItem = item
		}
	}

	payload := fortnite.UnlockRewardNodePayload{
		NodeID: rewardNodeID,
		RewardGraphID: rewardGraphID,
		RewardCFG: "",
	}

	unlockRes, err := client.UnlockRewardNode(payload)
	if err != nil {
		return fmt.Errorf("unable to UnlockRewardNode (athena): %s", err)
	}

	println(unlockRes)
	return nil
}