package commands

import (
	"blast/api"
	"blast/db"
	"encoding/json"
	"fmt"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	djson "github.com/disgoorg/json"
)

var daily = discord.SlashCommandCreate{
	Name:        "daily",
	Description: "Claims your Save the World daily login reward.",
}

var Daily = Command{
	Handler: func(event *handler.CommandEvent, blast api.EpicClient, user db.UserEntry, credentials api.UserCredentialsResponse, data discord.SlashCommandInteractionData) error {
		campaignProfile, err := blast.ProfileOperationStr(credentials, "ClaimLoginReward", "campaign", "{}")
		if err != nil {
			return err
		}

		defer campaignProfile.Close()

		var profile api.CampaignProfile
		err = json.NewDecoder(campaignProfile).Decode(&profile)
		if err != nil {
			return err
		}

		day := profile.ProfileChanges[0].Profile.Stats.Attributes.DailyRewards.TotalDaysLoggedIn

		account, err := blast.FetchAccountInformation(credentials)
		if err != nil {
			return err
		}

		avatarURL, err := blast.FetchAvatarURL(credentials)
		if err != nil {
			return err
		}

		todaysReward := RewardGraph[day%336].Item

		fields := []discord.EmbedField{}

		dayOffset := day % 7
		weeklyRewardsStr := ""

		for i := 0; i < 7; i++ {
			dayIndex := day - dayOffset + i + 1

			if dayIndex <= day {
				weeklyRewardsStr += fmt.Sprintf("`%d:` ~~%s~~", dayIndex, RewardGraph[dayIndex%336].Item)
			} else {
				weeklyRewardsStr += fmt.Sprintf("`%d:` %s", dayIndex, RewardGraph[dayIndex%336].Item)
			}

			if i != 6 {
				weeklyRewardsStr += "\n"
			}
		}

		fields = append(fields, discord.EmbedField{
			Name:   "This Week's Rewards",
			Value:  weeklyRewardsStr,
			Inline: djson.Ptr(true),
		})

		if len(profile.Notifications) > 0 {
			if len(profile.Notifications[0].Items) > 0 {
				fields = append(fields, discord.EmbedField{
					Name:   "Today's Reward Claimed",
					Value:  fmt.Sprintf("`%d:` %s", day, todaysReward),
					Inline: djson.Ptr(true),
				})
			}
		}

		_, err = event.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
			SetEmbeds(discord.NewEmbedBuilder().
				SetFields(fields...).
				SetAuthorIcon(avatarURL).
				SetColor(0xFB5A32).
				SetTimestamp(time.Now()).
				SetAuthorNamef(account.DisplayName).
				SetTitle("Daily Rewards").
				Build(),
			).
			Build(),
		)
		return err
	},
	LoginRequired:     true,
	EphemeralResponse: false,
}

type RewardGraphItem struct {
	Day      int
	Item     string
	EmojiStr string
}

var RewardGraph = []RewardGraphItem{
	{
		Day:      1,
		Item:     "1000x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      2,
		Item:     "300x Hero XP",
		EmojiStr: "",
	},
	{
		Day:      3,
		Item:     "2x  Mini Llamas",
		EmojiStr: "",
	},
	{
		Day:      4,
		Item:     "Rare Hero",
		EmojiStr: "",
	},
	{
		Day:      5,
		Item:     "Rare Ranged Weapon",
		EmojiStr: "",
	},
	{
		Day:      6,
		Item:     "Rare Melee Weapon",
		EmojiStr: "",
	},
	{
		Day:      7,
		Item:     "5x Pure Drops of Rain",
		EmojiStr: "",
	},
	{
		Day:      8,
		Item:     "Rare Freedom's Herald Pistol",
		EmojiStr: "",
	},
	{
		Day:      9,
		Item:     "Rare Trap",
		EmojiStr: "",
	},
	{
		Day:      10,
		Item:     "Rare Defender",
		EmojiStr: "",
	},
	{
		Day:      11,
		Item:     "Epic Hero",
		EmojiStr: "",
	},
	{
		Day:      12,
		Item:     "50x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      13,
		Item:     "Epic Lead Survivor",
		EmojiStr: "",
	},
	{
		Day:      14,
		Item:     "3x XP Boosts",
		EmojiStr: "",
	},
	{
		Day:      15,
		Item:     "Epic Survivor",
		EmojiStr: "",
	},
	{
		Day:      16,
		Item:     "Mini Llama",
		EmojiStr: "",
	},
	{
		Day:      17,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      18,
		Item:     "5x Pure Drops of Rain",
		EmojiStr: "",
	},
	{
		Day:      19,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      20,
		Item:     "Rare Trap",
		EmojiStr: "",
	},
	{
		Day:      21,
		Item:     "5x Lightning in a Bottle",
		EmojiStr: "",
	},
	{
		Day:      22,
		Item:     "Legendary Hero",
		EmojiStr: "",
	},
	{
		Day:      23,
		Item:     "5x Pure Drops of Rain",
		EmojiStr: "",
	},
	{
		Day:      24,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      25,
		Item:     "2x Mini Llamas",
		EmojiStr: "",
	},
	{
		Day:      26,
		Item:     "Rare Melee Weapon",
		EmojiStr: "",
	},
	{
		Day:      27,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      28,
		Item:     "5x Eye of the Storm",
		EmojiStr: "",
	},
	{
		Day:      29,
		Item:     "300x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      30,
		Item:     "2x Mini Llamas",
		EmojiStr: "",
	},
	{
		Day:      31,
		Item:     "10x Pure Drops of Rain ",
		EmojiStr: "",
	},
	{
		Day:      32,
		Item:     "Rare Hero",
		EmojiStr: "",
	},
	{
		Day:      33,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      34,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      35,
		Item:     "3x XP Boosts",
		EmojiStr: "",
	},
	{
		Day:      36,
		Item:     "150x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      37,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      38,
		Item:     "10x Pure Drops of Rain ",
		EmojiStr: "",
	},
	{
		Day:      39,
		Item:     "Rare Ranged Weapon",
		EmojiStr: "",
	},
	{
		Day:      40,
		Item:     "5x Lightning in a Bottle",
		EmojiStr: "",
	},
	{
		Day:      41,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      42,
		Item:     "3x Teammate XP Boosts",
		EmojiStr: "",
	},
	{
		Day:      43,
		Item:     "Epic Lead Survivor",
		EmojiStr: "",
	},
	{
		Day:      44,
		Item:     "2x Mini Llamas",
		EmojiStr: "",
	},
	{
		Day:      45,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      46,
		Item:     "5x Pure Drops of Rain",
		EmojiStr: "",
	},
	{
		Day:      47,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      48,
		Item:     "Rare Trap",
		EmojiStr: "",
	},
	{
		Day:      49,
		Item:     "5x Lightning in a Bottle",
		EmojiStr: "",
	},
	{
		Day:      50,
		Item:     "150x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      51,
		Item:     "5x Pure Drops of Rain",
		EmojiStr: "",
	},
	{
		Day:      52,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      53,
		Item:     "2x Mini Llamas",
		EmojiStr: "",
	},
	{
		Day:      54,
		Item:     "Rare Melee Weapon",
		EmojiStr: "",
	},
	{
		Day:      55,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      56,
		Item:     "5x Eye of the Storm",
		EmojiStr: "",
	},
	{
		Day:      57,
		Item:     "300x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      58,
		Item:     "2x Mini Llamas",
		EmojiStr: "",
	},
	{
		Day:      59,
		Item:     "10x Pure Drops of Rain ",
		EmojiStr: "",
	},
	{
		Day:      60,
		Item:     "Rare Hero",
		EmojiStr: "",
	},
	{
		Day:      61,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      62,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      63,
		Item:     "3x XP Boosts",
		EmojiStr: "",
	},
	{
		Day:      64,
		Item:     "Epic Survivor",
		EmojiStr: "",
	},
	{
		Day:      65,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      66,
		Item:     "10x Pure Drops of Rain ",
		EmojiStr: "",
	},
	{
		Day:      67,
		Item:     "Rare Ranged Weapon",
		EmojiStr: "",
	},
	{
		Day:      68,
		Item:     "5x Lightning in a Bottle",
		EmojiStr: "",
	},
	{
		Day:      69,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      70,
		Item:     "3x Teammate XP Boosts",
		EmojiStr: "",
	},
	{
		Day:      71,
		Item:     "150x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      72,
		Item:     "2x Mini Llamas",
		EmojiStr: "",
	},
	{
		Day:      73,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      74,
		Item:     "5x Pure Drops of Rain",
		EmojiStr: "",
	},
	{
		Day:      75,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      76,
		Item:     "Rare Trap",
		EmojiStr: "",
	},
	{
		Day:      77,
		Item:     "5x Lightning in a Bottle",
		EmojiStr: "",
	},
	{
		Day:      78,
		Item:     "150x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      79,
		Item:     "5x Pure Drops of Rain",
		EmojiStr: "",
	},
	{
		Day:      80,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      81,
		Item:     "2x Mini Llamas",
		EmojiStr: "",
	},
	{
		Day:      82,
		Item:     "Epic Survivor",
		EmojiStr: "",
	},
	{
		Day:      83,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      84,
		Item:     "5x Eye of the Storm",
		EmojiStr: "",
	},
	{
		Day:      85,
		Item:     "300x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      86,
		Item:     "2x Mini Llamas",
		EmojiStr: "",
	},
	{
		Day:      87,
		Item:     "10x Pure Drops of Rain ",
		EmojiStr: "",
	},
	{
		Day:      88,
		Item:     "Rare Hero",
		EmojiStr: "",
	},
	{
		Day:      89,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      90,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      91,
		Item:     "3x XP Boosts",
		EmojiStr: "",
	},
	{
		Day:      92,
		Item:     "150x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      93,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      94,
		Item:     "10x Pure Drops of Rain ",
		EmojiStr: "",
	},
	{
		Day:      95,
		Item:     "Rare Ranged Weapon",
		EmojiStr: "",
	},
	{
		Day:      96,
		Item:     "5x Lightning in a Bottle",
		EmojiStr: "",
	},
	{
		Day:      97,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      98,
		Item:     "3x Teammate XP Boosts",
		EmojiStr: "",
	},
	{
		Day:      99,
		Item:     "150x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      100,
		Item:     "2x Mini Llamas",
		EmojiStr: "",
	},
	{
		Day:      101,
		Item:     "Legendary Hero",
		EmojiStr: "",
	},
	{
		Day:      102,
		Item:     "5x Pure Drops of Rain",
		EmojiStr: "",
	},
	{
		Day:      103,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      104,
		Item:     "Rare Trap",
		EmojiStr: "",
	},
	{
		Day:      105,
		Item:     "5x Lightning in a Bottle",
		EmojiStr: "",
	},
	{
		Day:      106,
		Item:     "150x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      107,
		Item:     "5x Pure Drops of Rain",
		EmojiStr: "",
	},
	{
		Day:      108,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      109,
		Item:     "2x Mini Llamas",
		EmojiStr: "",
	},
	{
		Day:      110,
		Item:     "Rare Melee Weapon",
		EmojiStr: "",
	},
	{
		Day:      111,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      112,
		Item:     "5x Eye of the Storm",
		EmojiStr: "",
	},
	{
		Day:      113,
		Item:     "800xV-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      114,
		Item:     "2x Mini Llamas",
		EmojiStr: "",
	},
	{
		Day:      115,
		Item:     "10x Pure Drops of Rain ",
		EmojiStr: "",
	},
	{
		Day:      116,
		Item:     "Rare Hero",
		EmojiStr: "",
	},
	{
		Day:      117,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      118,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      119,
		Item:     "3x XP Boosts",
		EmojiStr: "",
	},
	{
		Day:      120,
		Item:     "150x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      121,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      122,
		Item:     "10x Pure Drops of Rain ",
		EmojiStr: "",
	},
	{
		Day:      123,
		Item:     "Rare Ranged Weapon",
		EmojiStr: "",
	},
	{
		Day:      124,
		Item:     "5x Lightning in a Bottle",
		EmojiStr: "",
	},
	{
		Day:      125,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      126,
		Item:     "3x Teammate XP Boosts",
		EmojiStr: "",
	},
	{
		Day:      127,
		Item:     "150x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      128,
		Item:     "2x Mini Llamas",
		EmojiStr: "",
	},
	{
		Day:      129,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      130,
		Item:     "5x Pure Drops of Rain",
		EmojiStr: "",
	},
	{
		Day:      131,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      132,
		Item:     "Rare Trap",
		EmojiStr: "",
	},
	{
		Day:      133,
		Item:     "5x Lightning in a Bottle",
		EmojiStr: "",
	},
	{
		Day:      134,
		Item:     "150x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      135,
		Item:     "5x Pure Drops of Rain",
		EmojiStr: "",
	},
	{
		Day:      136,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      137,
		Item:     "2x Mini Llamas",
		EmojiStr: "",
	},
	{
		Day:      138,
		Item:     "Rare Melee Weapon",
		EmojiStr: "",
	},
	{
		Day:      139,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      140,
		Item:     "5x Eye of the Storm",
		EmojiStr: "",
	},
	{
		Day:      141,
		Item:     "300x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      142,
		Item:     "2x Mini Llamas",
		EmojiStr: "",
	},
	{
		Day:      143,
		Item:     "10x Pure Drops of Rain ",
		EmojiStr: "",
	},
	{
		Day:      144,
		Item:     "Rare Hero",
		EmojiStr: "",
	},
	{
		Day:      145,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      146,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      147,
		Item:     "3x XP Boosts",
		EmojiStr: "",
	},
	{
		Day:      148,
		Item:     "150x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      149,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      150,
		Item:     "10x Pure Drops of Rain ",
		EmojiStr: "",
	},
	{
		Day:      151,
		Item:     "Rare Ranged Weapon",
		EmojiStr: "",
	},
	{
		Day:      152,
		Item:     "5x Lightning in a Bottle",
		EmojiStr: "",
	},
	{
		Day:      153,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      154,
		Item:     "3x Teammate XP Boosts",
		EmojiStr: "",
	},
	{
		Day:      155,
		Item:     "150x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      156,
		Item:     "2x Mini Llamas",
		EmojiStr: "",
	},
	{
		Day:      157,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      158,
		Item:     "5x Pure Drops of Rain",
		EmojiStr: "",
	},
	{
		Day:      159,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      160,
		Item:     "Rare Trap",
		EmojiStr: "",
	},
	{
		Day:      161,
		Item:     "5x Lightning in a Bottle",
		EmojiStr: "",
	},
	{
		Day:      162,
		Item:     "150x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      163,
		Item:     "5x Pure Drops of Rain",
		EmojiStr: "",
	},
	{
		Day:      164,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      165,
		Item:     "2x Mini Llamas",
		EmojiStr: "",
	},
	{
		Day:      166,
		Item:     "Rare Melee Weapon",
		EmojiStr: "",
	},
	{
		Day:      167,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      168,
		Item:     "5x Eye of the Storm",
		EmojiStr: "",
	},
	{
		Day:      169,
		Item:     "300x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      170,
		Item:     "2x Mini Llamas",
		EmojiStr: "",
	},
	{
		Day:      171,
		Item:     "10x Pure Drops of Rain ",
		EmojiStr: "",
	},
	{
		Day:      172,
		Item:     "Rare Hero",
		EmojiStr: "",
	},
	{
		Day:      173,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      174,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      175,
		Item:     "3x XP Boosts",
		EmojiStr: "",
	},
	{
		Day:      176,
		Item:     "150x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      177,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      178,
		Item:     "10x Pure Drops of Rain ",
		EmojiStr: "",
	},
	{
		Day:      179,
		Item:     "Rare Ranged Weapon",
		EmojiStr: "",
	},
	{
		Day:      180,
		Item:     "5x Lightning in a Bottle",
		EmojiStr: "",
	},
	{
		Day:      181,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      182,
		Item:     "3x Teammate XP Boosts",
		EmojiStr: "",
	},
	{
		Day:      183,
		Item:     "150x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      184,
		Item:     "2x Mini Llamas",
		EmojiStr: "",
	},
	{
		Day:      185,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      186,
		Item:     "15x Pure Drops of Rain",
		EmojiStr: "",
	},
	{
		Day:      187,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      188,
		Item:     "Rare Trap",
		EmojiStr: "",
	},
	{
		Day:      189,
		Item:     "5x Lightning in a Bottle",
		EmojiStr: "",
	},
	{
		Day:      190,
		Item:     "150x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      191,
		Item:     "5x Pure Drops of Rain",
		EmojiStr: "",
	},
	{
		Day:      192,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      193,
		Item:     "2x Mini Llamas",
		EmojiStr: "",
	},
	{
		Day:      194,
		Item:     "Rare Melee Weapon",
		EmojiStr: "",
	},
	{
		Day:      195,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      196,
		Item:     "5x Eye of the Storm",
		EmojiStr: "",
	},
	{
		Day:      197,
		Item:     "300x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      198,
		Item:     "2x Mini Llamas",
		EmojiStr: "",
	},
	{
		Day:      199,
		Item:     "10x Pure Drops of Rain ",
		EmojiStr: "",
	},
	{
		Day:      200,
		Item:     "Rare Hero",
		EmojiStr: "",
	},
	{
		Day:      201,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      202,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      203,
		Item:     "3x XP Boosts",
		EmojiStr: "",
	},
	{
		Day:      204,
		Item:     "150x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      205,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      206,
		Item:     "10x Pure Drops of Rain ",
		EmojiStr: "",
	},
	{
		Day:      207,
		Item:     "Rare Ranged Weapon",
		EmojiStr: "",
	},
	{
		Day:      208,
		Item:     "5x Lightning in a Bottle",
		EmojiStr: "",
	},
	{
		Day:      209,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      210,
		Item:     "3x Teammate XP Boosts",
		EmojiStr: "",
	},
	{
		Day:      211,
		Item:     "150x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      212,
		Item:     "2x Mini Llamas",
		EmojiStr: "",
	},
	{
		Day:      213,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      214,
		Item:     "5x Pure Drops of Rain",
		EmojiStr: "",
	},
	{
		Day:      215,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      216,
		Item:     "Rare Trap",
		EmojiStr: "",
	},
	{
		Day:      217,
		Item:     "5x Lightning in a Bottle",
		EmojiStr: "",
	},
	{
		Day:      218,
		Item:     "150x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      219,
		Item:     "5x Pure Drops of Rain",
		EmojiStr: "",
	},
	{
		Day:      220,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      221,
		Item:     "2x Mini Llamas",
		EmojiStr: "",
	},
	{
		Day:      222,
		Item:     "Rare Melee Weapon",
		EmojiStr: "",
	},
	{
		Day:      223,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      224,
		Item:     "5x Eye of the Storm",
		EmojiStr: "",
	},
	{
		Day:      225,
		Item:     "800xV-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      226,
		Item:     "2x Mini Llamas",
		EmojiStr: "",
	},
	{
		Day:      227,
		Item:     "10x Pure Drops of Rain ",
		EmojiStr: "",
	},
	{
		Day:      228,
		Item:     "Rare Hero",
		EmojiStr: "",
	},
	{
		Day:      229,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      230,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      231,
		Item:     "3x XP Boosts",
		EmojiStr: "",
	},
	{
		Day:      232,
		Item:     "150x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      233,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      234,
		Item:     "10x Pure Drops of Rain ",
		EmojiStr: "",
	},
	{
		Day:      235,
		Item:     "Rare Ranged Weapon",
		EmojiStr: "",
	},
	{
		Day:      236,
		Item:     "5x Lightning in a Bottle",
		EmojiStr: "",
	},
	{
		Day:      237,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      238,
		Item:     "3x Teammate XP Boosts",
		EmojiStr: "",
	},
	{
		Day:      239,
		Item:     "150x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      240,
		Item:     "2x Mini Llamas",
		EmojiStr: "",
	},
	{
		Day:      241,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      242,
		Item:     "5x Pure Drops of Rain",
		EmojiStr: "",
	},
	{
		Day:      243,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      244,
		Item:     "Rare Trap",
		EmojiStr: "",
	},
	{
		Day:      245,
		Item:     "5x Lightning in a Bottle",
		EmojiStr: "",
	},
	{
		Day:      246,
		Item:     "150x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      247,
		Item:     "5x Pure Drops of Rain",
		EmojiStr: "",
	},
	{
		Day:      248,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      249,
		Item:     "2x Mini Llamas",
		EmojiStr: "",
	},
	{
		Day:      250,
		Item:     "Rare Melee Weapon",
		EmojiStr: "",
	},
	{
		Day:      251,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      252,
		Item:     "5x Eye of the Storm",
		EmojiStr: "",
	},
	{
		Day:      253,
		Item:     "300x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      254,
		Item:     "2x Mini Llamas",
		EmojiStr: "",
	},
	{
		Day:      255,
		Item:     "10x Pure Drops of Rain ",
		EmojiStr: "",
	},
	{
		Day:      256,
		Item:     "Rare Hero",
		EmojiStr: "",
	},
	{
		Day:      257,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      258,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      259,
		Item:     "3x XP Boosts",
		EmojiStr: "",
	},
	{
		Day:      260,
		Item:     "150x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      261,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      262,
		Item:     "10x Pure Drops of Rain ",
		EmojiStr: "",
	},
	{
		Day:      263,
		Item:     "Rare Ranged Weapon",
		EmojiStr: "",
	},
	{
		Day:      264,
		Item:     "5x Lightning in a Bottle",
		EmojiStr: "",
	},
	{
		Day:      265,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      266,
		Item:     "3x Teammate XP Boosts",
		EmojiStr: "",
	},
	{
		Day:      267,
		Item:     "150x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      268,
		Item:     "2x Mini Llamas",
		EmojiStr: "",
	},
	{
		Day:      269,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      270,
		Item:     "5x Pure Drops of Rain",
		EmojiStr: "",
	},
	{
		Day:      271,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      272,
		Item:     "Rare Trap",
		EmojiStr: "",
	},
	{
		Day:      273,
		Item:     "5x Lightning in a Bottle",
		EmojiStr: "",
	},
	{
		Day:      274,
		Item:     "150x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      275,
		Item:     "5x Pure Drops of Rain",
		EmojiStr: "",
	},
	{
		Day:      276,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      277,
		Item:     "2x Mini Llamas",
		EmojiStr: "",
	},
	{
		Day:      278,
		Item:     "Rare Melee Weapon",
		EmojiStr: "",
	},
	{
		Day:      279,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      280,
		Item:     "5x Eye of the Storm",
		EmojiStr: "",
	},
	{
		Day:      281,
		Item:     "300x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      282,
		Item:     "2x Mini Llamas",
		EmojiStr: "",
	},
	{
		Day:      283,
		Item:     "10x Pure Drops of Rain ",
		EmojiStr: "",
	},
	{
		Day:      284,
		Item:     "Rare Hero",
		EmojiStr: "",
	},
	{
		Day:      285,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      286,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      287,
		Item:     "3x XP Boosts",
		EmojiStr: "",
	},
	{
		Day:      288,
		Item:     "150x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      289,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      290,
		Item:     "10x Pure Drops of Rain ",
		EmojiStr: "",
	},
	{
		Day:      291,
		Item:     "Rare Ranged Weapon",
		EmojiStr: "",
	},
	{
		Day:      292,
		Item:     "5x Lightning in a Bottle",
		EmojiStr: "",
	},
	{
		Day:      293,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      294,
		Item:     "3x Teammate XP Boosts",
		EmojiStr: "",
	},
	{
		Day:      295,
		Item:     "150x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      296,
		Item:     "2x Mini Llamas",
		EmojiStr: "",
	},
	{
		Day:      297,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      298,
		Item:     "5x Pure Drops of Rain",
		EmojiStr: "",
	},
	{
		Day:      299,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      300,
		Item:     "Rare Trap",
		EmojiStr: "",
	},
	{
		Day:      301,
		Item:     "5x Lightning in a Bottle",
		EmojiStr: "",
	},
	{
		Day:      302,
		Item:     "150x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      303,
		Item:     "5x Pure Drops of Rain",
		EmojiStr: "",
	},
	{
		Day:      304,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      305,
		Item:     "2x Mini Llamas",
		EmojiStr: "",
	},
	{
		Day:      306,
		Item:     "Rare Melee Weapon",
		EmojiStr: "",
	},
	{
		Day:      307,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      308,
		Item:     "5x Eye of the Storm",
		EmojiStr: "",
	},
	{
		Day:      309,
		Item:     "300x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      310,
		Item:     "2x Mini Llamas",
		EmojiStr: "",
	},
	{
		Day:      311,
		Item:     "10x Pure Drops of Rain ",
		EmojiStr: "",
	},
	{
		Day:      312,
		Item:     "Rare Hero",
		EmojiStr: "",
	},
	{
		Day:      313,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      314,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      315,
		Item:     "3x XP Boosts",
		EmojiStr: "",
	},
	{
		Day:      316,
		Item:     "150x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      317,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      318,
		Item:     "10x Pure Drops of Rain ",
		EmojiStr: "",
	},
	{
		Day:      319,
		Item:     "Rare Ranged Weapon",
		EmojiStr: "",
	},
	{
		Day:      320,
		Item:     "5x Lightning in a Bottle",
		EmojiStr: "",
	},
	{
		Day:      321,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      322,
		Item:     "3x Teammate XP Boosts",
		EmojiStr: "",
	},
	{
		Day:      323,
		Item:     "150x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      324,
		Item:     "2x Mini Llamas",
		EmojiStr: "",
	},
	{
		Day:      325,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      326,
		Item:     "5x Pure Drops of Rain",
		EmojiStr: "",
	},
	{
		Day:      327,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      328,
		Item:     "Rare Trap",
		EmojiStr: "",
	},
	{
		Day:      329,
		Item:     "5x Lightning in a Bottle",
		EmojiStr: "",
	},
	{
		Day:      330,
		Item:     "150x V-Bucks & X-Ray Tickets",
		EmojiStr: "",
	},
	{
		Day:      331,
		Item:     "5x Pure Drops of Rain",
		EmojiStr: "",
	},
	{
		Day:      332,
		Item:     "20x Gold",
		EmojiStr: "",
	},
	{
		Day:      333,
		Item:     "2x Mini Llamas",
		EmojiStr: "",
	},
	{
		Day:      334,
		Item:     "Rare Melee Weapon",
		EmojiStr: "",
	},
	{
		Day:      335,
		Item:     "Upgrade Llama",
		EmojiStr: "",
	},
	{
		Day:      336,
		Item:     "5x Eye of the Storm",
		EmojiStr: "",
	},
}
