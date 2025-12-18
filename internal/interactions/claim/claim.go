package claim

import (
	"encoding/json"
	"fmt"
	"math"
	"slices"
	"strings"
	"time"

	"github.com/8h9x/BlastBot/internal/manager/sessions"
	"github.com/8h9x/fortgo/fortnite"
	"github.com/8h9x/fortgo/request"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

const (
	rewardGraph = "AthenaRewardGraph:s39_winterfest"
	specialRewardNode = "ERG.Node.B.1"
)

type PseduoRewardGraphItemAttrributes struct {
	RewardNodesClaimed []string `json:"reward_nodes_claimed"`
	RewardKeys         []struct {
		UnlockKeysUsed int `json:"unlock_keys_used"`
	} `json:"reward_keys"`
}

type PseduoRewardGraphItem struct {
	TemplateID string `json:"templateId"`
	Attributes PseduoRewardGraphItemAttrributes `json:"attributes"`
	Quantity   int  `json:"quantity"`
}

var Definition = discord.SlashCommandCreate{
	Name:        "claim",
	Description: "Claims the {x} outfit winterfest present early.)",
}

func Handler(event *handler.CommandEvent) error {
	discordId := event.User().ID

	session, err := sessions.GetSessionForUser(discordId)
	if err != nil {
		return fmt.Errorf("unable to get session for user (%s): %s", discordId, err)
	}

	res, err := session.QueryProfile("athena")
	if err != nil {
		return fmt.Errorf("unable to QueryProfile (athena): %s", err)
	}

	profileRes, err := request.ParseResponse[fortnite.Profile[fortnite.AthenaProfileStats, []any]](res)
	if err != nil {
		return fmt.Errorf("unable to parse QueryProfile response: %s", err)
	}

	items := profileRes.Data.ProfileChanges[0].Profile.Items

	var rewardGraphID string
	var rewardGraphItem PseduoRewardGraphItem

	for key, value := range items {
		var item PseduoRewardGraphItem
		if err = json.Unmarshal(value, &item); err != nil {
			// not a skin; (you should probably add an additional check to ensure that it isn't some other type of error occurring); TODO: abstract this to a helper function that properly error checks and returns an empty state of the type passed if the type of data doesnt match
			continue
		}

		if strings.HasPrefix(item.TemplateID, rewardGraph) {
			rewardGraphID = key
			rewardGraphItem = item
		}
	}

	unlockEpoch := time.Date(2025, 12, 18, 14, 0, 0, 0, time.UTC).UnixMilli()
	now := time.Now().UnixMilli()
	day := int64(24 * time.Hour / time.Millisecond)
	expectedKeysAmount := int(math.Round(float64(now-unlockEpoch) / float64(day)))

	if slices.Contains(rewardGraphItem.Attributes.RewardNodesClaimed, specialRewardNode) {
		embed := discord.NewEmbedBuilder().
			SetColor(0xFF3333).
			SetTimestamp(time.Now()).
			SetTitle("Present Already Opened!").
			SetDescriptionf("It looks like you already have the '%s' outfit present opened.\n\nIf the item isnt showing in your locker, try restarting the game.", "{x}")

		err = event.CreateMessage(discord.NewMessageCreateBuilder().
			SetEmbeds(embed.
				Build(),
			).
			Build(),
		)
		if err != nil {
			return err
		}
	} else if len(rewardGraphItem.Attributes.RewardKeys) == 0 || rewardGraphItem.Attributes.RewardKeys[0].UnlockKeysUsed >= expectedKeysAmount {
		embed := discord.NewEmbedBuilder().
			SetColor(0xFF3333).
			SetTimestamp(time.Now()).
			SetTitle("Present Not Available Yet!").
			SetDescriptionf("It looks like you don't have enough gifts to open the '%s' outfit present yet.\n\nYou should be able to open a new gift tomorrow at 9am EST.", "{x}")

		err = event.CreateMessage(discord.NewMessageCreateBuilder().
			SetEmbeds(embed.
				Build(),
			).
			Build(),
		)
		if err != nil {
			return err
		}
	}

	payload := fortnite.UnlockRewardNodePayload{
		NodeID: specialRewardNode,
		RewardGraphID: rewardGraphID,
		RewardCFG: "",
	}

	unlockRes, err := session.UnlockRewardNode(payload)
	if err != nil {
		return fmt.Errorf("unable to UnlockRewardNode (athena): %s", err)
	}

	println(unlockRes)

	embed := discord.NewEmbedBuilder().
		SetColor(0x00C059).
		SetTimestamp(time.Now()).
		SetTitlef("Opened the '%s' present.", "{x}").
		SetDescriptionf("The '%s' winterfest gift has been claimed. Please restart your game for the reward screen to appear.", "{x}")

	err = event.CreateMessage(discord.NewMessageCreateBuilder().
		SetEmbeds(embed.
			Build(),
		).
		Build(),
	)
	if err != nil {
		return err
	}

	return nil
}