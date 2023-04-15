package commands

import (
	"blast/api"
	"blast/db"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var account = discord.SlashCommandCreate{
	Name:        "account",
	Description: "Account related commands.",

	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionSubCommand{
			Name:        "info",
			Description: "Retrieves general account information.",
		},
	},
}

var AccountInfo = Command{
	Handler: func(event *handler.CommandEvent, blast api.EpicClient, user db.UserEntry, credentials api.UserCredentialsResponse, data discord.SlashCommandInteractionData) error {
		account, err := blast.FetchAccountInformation(credentials)
		if err != nil {
			return err
		}

		brInfo, err := blast.FetchAccountBRInfo(credentials)
		if err != nil {
			return err
		}

		athenaProfile, err := blast.ProfileOperationStr(credentials, "QueryProfile", "athena", "{}")
		if err != nil {
			return err
		}

		defer athenaProfile.Close()

		var profile api.AthenaProfile[api.AthenaProfileLockerItem]
		err = json.NewDecoder(athenaProfile).Decode(&profile)
		if err != nil {
			return err
		}

		attributes := profile.ProfileChanges[0].Profile.Stats.Attributes
		lastAppliedLoadout := attributes.LastAppliedLoadout
		locker := profile.ProfileChanges[0].Profile.Items[lastAppliedLoadout]

		bpEmoji := "<:free_pass:1096479702417416243>"
		if attributes.BookPurchased {
			bpEmoji = "<:battlepass:1096473607447777383>"
		}

		skinCount := 0
		for _, slot := range profile.ProfileChanges[0].Profile.Items {
			if strings.HasPrefix(slot.TemplateID, "AthenaCharacter") {
				skinCount++
			}
		}

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

		// partyMetaUpdate := api.PartyMetaUpdate{
		// 	Delete:   []string{},
		// 	Revision: 4,
		// 	Update: api.PartyMetaUpdateData{
		// 		DefaultAthenaCosmeticLoadoutJ: api.PartyMetaUpdateDefaultAthenaCosmeticLoadoutJ{
		// 			AthenaCosmeticLoadout: api.PartyMetaUpdateDefaultAthenaCosmeticLoadout{
		// 				CharacterDef:  "/Game/Athena/Items/Cosmetics/Characters/CID_342_Athena_Commando_M_StreetRacerMetallic.CID_342_Athena_Commando_M_StreetRacerMetallic",
		// 				CharacterEKey: "",
		// 				BackpackDef:   "/Game/Athena/Items/Cosmetics/Backpacks/BID_610_ElasticHologram.BID_610_ElasticHologram",
		// 				BackpackEKey:  "",
		// 				PickaxeDef:    "/Game/Athena/Items/Cosmetics/Pickaxes/Pickaxe_ID_035_Prismatic.Pickaxe_ID_035_Prismatic",
		// 				PickaxeEKey:   "",
		// 				ContrailDef:   "/Game/Athena/Items/Cosmetics/Contrails/Trails_ID_019_PSBurnout.Trails_ID_019_PSBurnout",
		// 				ContrailEKey:  "",
		// 				Scratchpad:    []string{},
		// 				CosmeticStats: []api.PartyMetaUpdateAthenaCosmeticStat{
		// 					{
		// 						StatName:  "TotalVictoryCrowns",
		// 						StatValue: 999,
		// 					},
		// 					{
		// 						StatName:  "TotalRoyalRoyales",
		// 						StatValue: 999,
		// 					},
		// 					{
		// 						StatName:  "HasCrown",
		// 						StatValue: 1,
		// 					},
		// 				},
		// 			},
		// 		},
		// 	},
		// }

		// json, err := json.Marshal(partyMetaUpdate)
		// if err != nil {
		// 	return err
		// }

		updateData := `{"overridden":{},"deleted":[],"revision":4,"updated":{"Default:AthenaCosmeticLoadout_j":{"AthenaCosmeticLoadout":{"characterDef":"/Game/Athena/Items/Cosmetics/Characters/CID_342_Athena_Commando_M_StreetRacerMetallic.CID_342_Athena_Commando_M_StreetRacerMetallic","characterEKey":"","backpackDef":"/Game/Athena/Items/Cosmetics/Backpacks/BID_610_ElasticHologram.BID_610_ElasticHologram","backpackEKey":"","pickaxeDef":"/Game/Athena/Items/Cosmetics/Pickaxes/Pickaxe_ID_035_Prismatic.Pickaxe_ID_035_Prismatic","pickaxeEKey":"","contrailDef":"/Game/Athena/Items/Cosmetics/Contrails/Trails_ID_019_PSBurnout.Trails_ID_019_PSBurnout","contrailEKey":"","scratchpad":[],"cosmeticStats":[{"statName":"TotalVictoryCrowns","statValue":999},{"statName":"TotalRoyalRoyales","statValue":999},{"statName":"HasCrown","statValue":1}]}}}}`

		party, err := blast.FetchParty(credentials)
		if err != nil {
			return err
		}

		res, err := blast.PartyMetaUpdate(credentials, party.Current[0].ID, updateData)
		if err != nil {
			return err
		}

		log.Println(res)

		embed := discord.NewEmbedBuilder().
			SetAuthorIconf("https://fortnite-api.com/images/cosmetics/br/%s/icon.png", strings.Replace(locker.Attributes.LockerSlotsData.Slots["Character"].Items[0], "AthenaCharacter:", "", -1)). // TODO set author icon to bot user avatar
			SetColor(0xFB5A32).
			SetTimestamp(time.Now()).
			SetAuthorNamef("%s | %s", account.DisplayName, credentials.AccountID).
			AddField("<:llama:1096476378121126000> Account Level", fmt.Sprint(attributes.AccountLevel), true).
			AddField(fmt.Sprintf("%s Battlepass Level", bpEmoji), fmt.Sprint(attributes.Level), true).
			AddField("<:battle_star:1096473613504368640> Battlestars", fmt.Sprint(attributes.Battlestars), true).
			AddField("<:xp:1096486771820347453> Season XP", fmt.Sprint(attributes.XP), true).
			AddField("<:goldbars:1096469831034863637> Bars", fmt.Sprint(brInfo.Stash.Globalcash), true).
			AddField("<:victory:1096481674872770591> Lifetime Wins", fmt.Sprint(attributes.LifetimeWins), true).
			// AddField("<:victory_crown:1096481681575260314> Crown Wins", "coming soon", true).
			AddField("<:outfit:1096486172655636561> Skin Count", fmt.Sprint(skinCount), true).
			AddField("<:lock:1096491497987260578> MFA Reward Claimed", boolToEmoji(attributes.MFARewardClaimed), true).
			AddField(fmt.Sprintf(":flag_%s: Region", strings.ToLower(account.Country)), account.Country, true).
			Build()

		_, err = event.Client().Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.NewMessageUpdateBuilder().SetEmbeds(embed).Build())
		return err
	},
	LoginRequired:     true,
	EphemeralResponse: false,
}

func boolToEmoji(b bool) string {
	if b {
		return "<a:thumbs_up:1096490549218906193>"
	}

	return "<a:thumbs_down:1096490737715130388>"
}
