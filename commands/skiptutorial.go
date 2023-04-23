package commands

import (
	"blast/api"
	"blast/db"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var skipTutorial = discord.SlashCommandCreate{
	Name:        "skiptutorial",
	Description: "Skips the Save the World tutorial & grants the music pack (even if you do not own Save the World).",
}

var SkipTutorial = Command{
	Handler: func(event *handler.CommandEvent, blast api.EpicClient, user db.UserEntry, credentials api.UserCredentialsResponse, data discord.SlashCommandInteractionData) error {
		account, err := blast.FetchMyAccountInfo(credentials)
		if err != nil {
			return err
		}

		avatarURL, err := blast.FetchAvatarURL(credentials)
		if err != nil {
			return err
		}

		campaignProfile, err := blast.ProfileOperationStr(credentials, "QueryProfile", "campaign", "{}")
		if err != nil {
			return err
		}

		defer campaignProfile.Close()

		var profile api.CampaignProfile
		err = json.NewDecoder(campaignProfile).Decode(&profile)
		if err != nil {
			return err
		}

		for _, item := range profile.ProfileChanges[0].Profile.Items {
			if strings.ToLower(item.TemplateID) == "quest:homebaseonboarding" {
				if strings.ToLower(item.Attributes.QuestState) == "claimed" {
					embed := discord.NewEmbedBuilder().
						SetAuthorIcon(avatarURL).
						SetColor(0xFB5A32).
						SetTimestamp(time.Now()).
						SetAuthorName(account.DisplayName).
						SetDescription("You have already completed the tutorial.").
						Build()

					_, err = event.Client().Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.NewMessageUpdateBuilder().SetEmbeds(embed).Build())
					return err
				} else {
					break
				}
			}
		}

		_, err = blast.ProfileOperationStr(credentials, "SkipTutorial", "campaign", "{}")
		if err != nil {
			return err
		}

		updateQuestClientObjectives := UpdateQuestClientObjectives{
			Advance: []UpdateQuestClientObjectivesAdvance{
				{
					StatName:        "hbonboarding_watchsatellitecine",
					Count:           1,
					TimestampOffset: 0,
				},
				{
					StatName:        "hbonboarding_namehomebase",
					Count:           1,
					TimestampOffset: 0,
				},
			},
		}

		body, err := json.Marshal(updateQuestClientObjectives)
		if err != nil {
			return err
		}

		updatedCampaignProfile, err := blast.ProfileOperationStr(credentials, "UpdateQuestClientObjectives", "campaign", string(body))
		if err != nil {
			return err
		}

		defer updatedCampaignProfile.Close()

		var updatedProfile api.CampaignProfile
		err = json.NewDecoder(updatedCampaignProfile).Decode(&updatedProfile)
		if err != nil {
			return err
		}

		for _, item := range updatedProfile.ProfileChanges[0].Profile.Items {
			if strings.ToLower(item.TemplateID) == "quest:homebaseonboarding" {
				if strings.ToLower(item.Attributes.QuestState) == "claimed" {
					embed := discord.NewEmbedBuilder().
						SetAuthorIcon(avatarURL).
						SetColor(0xFB5A32).
						SetTimestamp(time.Now()).
						SetAuthorName(account.DisplayName).
						SetDescription("Successfully skipped the tutorial mission and granted your account the **Save the World** music pack.").
						Build()

					_, err = event.Client().Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.NewMessageUpdateBuilder().SetEmbeds(embed).Build())
					return err
				} else {
					break
				}
			}
		}

		return errors.New("an unknown error occured")
	},
	LoginRequired:     true,
	EphemeralResponse: false,
}

type UpdateQuestClientObjectives struct {
	Advance []UpdateQuestClientObjectivesAdvance `json:"advance"`
}

type UpdateQuestClientObjectivesAdvance struct {
	StatName        string `json:"statName"`        // The stat Name, e.g phoenixtheater_power_highlight
	Count           int    `json:"count"`           // The count to what u want to update it
	TimestampOffset int    `json:"timestampOffset"` // Leave it 0, this is related to the difference between Local and UTC Time
}

// {
//     "advance": [{
//         "statName": "phoenixtheater_power_highlight", // The stat Name, e.g phoenixtheater_power_highlight
//         "count": 1, // The count to what u want to update it
//         "timestampOffset": 0 // Leave it 0, this is related to the difference between Local and UTC Time
//     }]
// }
