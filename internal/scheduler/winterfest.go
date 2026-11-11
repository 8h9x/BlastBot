package scheduler

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/robfig/cron"

	"github.com/8h9x/fortgo"
	"github.com/8h9x/fortgo/fortnite"
	"github.com/8h9x/fortgo/mcp"
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

type WinterfestGiftService struct {
	cron                  cron.Cron
	RewardGraphTemplateID string
}

func (s *WinterfestGiftService) Register() error {
	return s.cron.AddFunc("0 30 * * * *", func() { s.Tick() })
}

func (s *WinterfestGiftService) Tick() {
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

	profileRes, err := request.ParseResponse[mcp.Profile](res)
	if err != nil {
		return fmt.Errorf("unable to parse ClientQuestLogin response: %s", err)
	}

	items := profileRes.Data.ProfileChanges[0].Profile.Items

	var rewardGraphID string

	for key, value := range items {
		if strings.HasPrefix(value.TemplateID, s.RewardGraphTemplateID) {
			rewardGraphID = key
		}
	}

	payload := fortnite.UnlockRewardNodePayload{
		NodeID:        rewardNodeID,
		RewardGraphID: rewardGraphID,
		RewardCFG:     "",
	}

	unlockRes, err := client.UnlockRewardNode(payload)
	if err != nil {
		return fmt.Errorf("unable to UnlockRewardNode (athena): %s", err)
	}

	println(unlockRes)
	return nil
}
