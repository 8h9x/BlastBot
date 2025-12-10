package commands

//import (
//	"context"
//	"encoding/json"
//	"fmt"
//
//	"github.com/8h9x/vinderman"
//	"github.com/8h9x/vinderman/consts"
//	"github.com/8h9x/vinderman/eos"
//	"github.com/8h9x/vinderman/request"
//	"go.mongodb.org/mongo-driver/v2/bson"
//
//	"strings"
//	"time"
//
//	"github.com/disgoorg/disgo/discord"
//	"github.com/disgoorg/disgo/handler"
//)
//
//var account = discord.SlashCommandCreate{
//	Name:        "account",
//	Description: "Account related commands.",
//	Options: []discord.ApplicationCommandOption{
//		discord.ApplicationCommandOptionSubCommand{
//			Name:        "info",
//			Description: "Retrieves general account information.",
//		},
//		discord.ApplicationCommandOptionSubCommand{
//			Name:        "switch",
//			Description: "Switches your selected epic games account.",
//		},
//		discord.ApplicationCommandOptionSubCommand{
//			Name:        "vbucks",
//			Description: "Check how many vbucks are on your account(s).",
//			Options: []discord.ApplicationCommandOption{
//				discord.ApplicationCommandOptionBool{
//					Name:        "bulk",
//					Description: "Whether to check all accounts or just the selected one.",
//				},
//			},
//		},
//	},
//}
//
//var AccountInfo = Command{
//	Handler: func(event *handler.CommandEvent, client *common.Client, user db.UserEntry, credentials vinderman.UserCredentials, data discord.SlashCommandInteractionData) error {
//		account, err := client.Epic.FetchMe(credentials)
//		if err != nil {
//			return err
//		}
//
//		brInfo, err := client.Epic.FetchBRInventory(credentials)
//		if err != nil {
//			return err
//		}
//
//		resp, err := client.Epic.QueryProfile(credentials, "athena")
//		if err != nil {
//			return err
//		}
//
//		res, err := request.ResponseParser[vinderman.Profile[vinderman.AthenaProfileStats, []interface{}]](resp)
//		if err != nil {
//			return err
//		}
//
//		profile := res.Body
//
//		attributes := profile.ProfileChanges[0].Profile.Stats.Attributes
//
//		bpEmoji := "<:free_pass:1096479702417416243>"
//		if attributes.BookPurchased {
//			bpEmoji = "<:battlepass:1096473607447777383>"
//		}
//
//		skinCount := 0
//		for _, item := range profile.ProfileChanges[0].Profile.Items {
//			var cosmetic vinderman.AthenaCosmeticItem
//			if err = json.Unmarshal(item, &cosmetic); err != nil {
//				// not a skin; (you should probably add an additional check to ensure that it isnt some other type of error occurring); TODO: abstract this to a helper function that properly error checks and returns an empty state of the type passed if the type of data doesnt match
//				continue
//			}
//
//			if strings.HasPrefix(cosmetic.TemplateID, "AthenaCharacter") {
//				skinCount++
//			}
//		}
//
//		avatarURL, err := client.Epic.FetchAvatarURL(credentials)
//		if err != nil {
//			return err
//		}
//
//		embed := discord.NewEmbedBuilder().
//			SetAuthorIcon(avatarURL).
//			SetColor(0xFB5A32).
//			SetTimestamp(time.Now()).
//			SetAuthorNamef("%s | %s", account.DisplayName, credentials.AccountID).
//			AddField("<:llama:1096476378121126000> Account Level", fmt.Sprint(attributes.AccountLevel), true).
//			AddField(fmt.Sprintf("%s Season Level", bpEmoji), fmt.Sprint(attributes.Level), true).
//			AddField("<:battle_star:1096473613504368640> Battlestars", fmt.Sprint(attributes.Battlestars), true).
//			AddField("<:xp:1096486771820347453> Season XP", fmt.Sprint(attributes.XP), true).
//			AddField("<:goldbars:1096469831034863637> Bars", fmt.Sprint(brInfo.Stash.Globalcash), true).
//			AddField("<:victory:1096481674872770591> Lifetime Wins", fmt.Sprint(attributes.LifetimeWins), true).
//			// AddField("<:victory_crown:1096481681575260314> Crown Wins", "coming soon", true).
//			AddField("<:outfit:1096486172655636561> Skin Count", fmt.Sprint(skinCount), true).
//			AddField("<:lock:1096491497987260578> MFA Reward Claimed", common.BoolToEmoji(attributes.MFARewardClaimed), true).
//			AddField(fmt.Sprintf(":flag_%s: Region", strings.ToLower(account.Country)), account.Country, true).
//			Build()
//
//		_, err = event.Client().Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.NewMessageUpdateBuilder().SetEmbeds(embed).Build())
//		return err
//	},
//	LoginRequired:     true,
//	EphemeralResponse: false,
//}
//
//var AccountSwitch = Command{
//	Handler: func(event *handler.CommandEvent, client *common.Client, user db.UserEntry, credentials vinderman.UserCredentials, data discord.SlashCommandInteractionData) error {
//		switchOptions := make([]discord.StringSelectMenuOption, 0)
//
//		for i, account := range user.Accounts {
//			refreshCredentials, err := client.Epic.RefreshTokenLogin(consts.FORTNITE_PC_CLIENT_ID, consts.FORTNITE_PC_CLIENT_SECRET, account.RefreshToken)
//			if err != nil {
//				if err.(*request.Error[eos.EpicErrorResponse]).Raw.ErrorCode == "errors.com.epicgames.account.auth_token.invalid_refresh_token" {
//					col := db.GetCollection("users")
//					_, err = col.UpdateOne(context.Background(), bson.M{"discordId": event.User().ID}, bson.M{"$pull": bson.M{"accounts": bson.M{"accountId": account.AccountID}}})
//					if err != nil {
//						return err
//					}
//
//					continue
//				}
//				return err
//			}
//
//			accountInfo, err := client.Epic.FetchMe(refreshCredentials)
//			if err != nil {
//				return err
//			}
//
//			switchOptions = append(switchOptions, discord.NewStringSelectMenuOption(accountInfo.DisplayName, fmt.Sprint(i)).WithDefault(i == user.SelectedAccount))
//		}
//
//		_, err := event.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
//			AddActionRow(discord.NewStringSelectMenu("switch_account_select", "Select the saved account to switch to", switchOptions...)).
//			AddActionRow(discord.NewDangerButton("Cancel", "cancel").WithEmoji(discord.ComponentEmoji{
//				Name: "x_",
//				ID:   1096630553689739385,
//			})).
//			Build(),
//		)
//		return err
//	},
//	LoginRequired:     true,
//	EphemeralResponse: false,
//}
//
//var AccountVbucks = Command{
//	Handler: func(event *handler.CommandEvent, client *common.Client, user db.UserEntry, credentials vinderman.UserCredentials, data discord.SlashCommandInteractionData) error {
//		if data.Bool("bulk") {
//			embed := discord.NewEmbedBuilder().SetColor(0xFB5A32).SetTimestamp(time.Now()).SetTitle("Saved Accounts V-Bucks")
//
//			bulkTotal := 0
//
//			for _, account := range user.Accounts {
//				refreshCredentials, err := client.Epic.RefreshTokenLogin(consts.FORTNITE_PC_CLIENT_ID, consts.FORTNITE_PC_CLIENT_SECRET, account.RefreshToken)
//				if err != nil {
//					if err.(*request.Error[eos.EpicErrorResponse]).Raw.ErrorCode == "errors.com.epicgames.account.auth_token.invalid_refresh_token" {
//						col := db.GetCollection("users")
//						_, err = col.UpdateOne(context.Background(), bson.M{"discordId": event.User().ID}, bson.M{"$pull": bson.M{"accounts": bson.M{"accountId": account.AccountID}}})
//						if err != nil {
//							return err
//						}
//
//						continue
//					}
//					return err
//				}
//
//				accountInfo, err := client.Epic.FetchMe(refreshCredentials)
//				if err != nil {
//					return err
//				}
//
//				resp, err := client.Epic.QueryProfile(credentials, "athena")
//				if err != nil {
//					return err
//				}
//
//				res, err := request.ResponseParser[vinderman.Profile[vinderman.CommonCoreProfileStats, []interface{}]](resp)
//				if err != nil {
//					return err
//				}
//
//				profile := res.Body
//
//				total := 0
//				list := ""
//
//				for _, entry := range profile.ProfileChanges[0].Profile.Items {
//					var item vinderman.CommonCoreMtxItem
//					if err = json.Unmarshal(entry, &item); err != nil {
//						// not a skin; (you should probably add an additional check to ensure that it isnt some other type of error occurring); TODO: abstract this to a helper function that properly error checks and returns an empty state of the type passed if the type of data doesnt match
//						continue
//					}
//
//					if strings.HasPrefix(item.TemplateID, "Currency:Mtx") {
//						total += item.Quantity
//						bulkTotal += item.Quantity
//
//						if len(user.Accounts) < 10 { // only show the extended list if there are less than 10 accounts
//							platform := VbuckPlatformsFriendly[item.Attributes.Platform]
//							if platform == "" {
//								platform = item.Attributes.Platform
//							}
//
//							list += fmt.Sprintf("<:vbuck:1099111673966628905> %d - *%s (%s)*\n", item.Quantity, platform, strings.Replace(item.TemplateID, "Currency:Mtx", "", 1))
//						}
//					}
//				}
//
//				if len(user.Accounts) < 10 {
//					embed.AddField(accountInfo.DisplayName, fmt.Sprintf("<:vbuck:1099111673966628905> **%d - Total** \n%s", total, list), true)
//				} else {
//					embed.AddField(accountInfo.DisplayName, fmt.Sprintf("<:vbuck:1099111673966628905> %d", total), true)
//				}
//			}
//
//			embed.SetFooterTextf("Total: %d V-Bucks", bulkTotal)
//
//			_, err := event.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().SetEmbeds(embed.Build()).Build())
//			return err
//		}
//
//		account, err := client.Epic.FetchMe(credentials)
//		if err != nil {
//			return err
//		}
//
//		avatarURL, err := client.Epic.FetchAvatarURL(credentials)
//		if err != nil {
//			return err
//		}
//
//		resp, err := client.Epic.QueryProfile(credentials, "athena")
//		if err != nil {
//			return err
//		}
//
//		res, err := request.ResponseParser[vinderman.Profile[vinderman.CommonCoreProfileStats, []interface{}]](resp)
//		if err != nil {
//			return err
//		}
//
//		profile := res.Body
//
//		total := 0
//		list := ""
//
//		for _, entry := range profile.ProfileChanges[0].Profile.Items {
//			var item vinderman.CommonCoreMtxItem
//			if err = json.Unmarshal(entry, &item); err != nil {
//				// not a skin; (you should probably add an additional check to ensure that it isn't some other type of error occurring); TODO: abstract this to a helper function that properly error checks and returns an empty state of the type passed if the type of data doesnt match
//				continue
//			}
//
//			if strings.HasPrefix(item.TemplateID, "Currency:Mtx") {
//				total += item.Quantity
//				platform := VbuckPlatformsFriendly[item.Attributes.Platform]
//				if platform == "" {
//					platform = item.Attributes.Platform
//				}
//
//				list += fmt.Sprintf("<:vbuck:1099111673966628905> **%d** - %s (%s)\n", item.Quantity, platform, strings.Replace(item.TemplateID, "Currency:Mtx", "", 1))
//			}
//		}
//
//		//for _, item := range profile.ProfileChanges[0].Profile.Items {
//		//	if strings.HasPrefix(item.TemplateID, "Currency:Mtx") {
//		//		total += item.Quantity
//		//		platform := VbuckPlatformsFriendly[item.Attributes.Platform]
//		//		if platform == "" {
//		//			platform = item.Attributes.Platform
//		//		}
//		//
//		//		list += fmt.Sprintf("<:vbuck:1099111673966628905> **%d** - %s (%s)\n", item.Quantity, platform, strings.Replace(item.TemplateID, "Currency:Mtx", "", 1))
//		//	}
//		//}
//
//		_, err = event.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
//			SetEmbeds(discord.NewEmbedBuilder().
//				SetTitle(fmt.Sprintf("<:vbuck:1099111673966628905> %d Total", total)).
//				SetColor(0xFB5A32).
//				SetAuthorNamef("%s's V-Bucks", account.DisplayName).
//				SetAuthorIcon(avatarURL).
//				SetDescription(list).
//				SetFooterTextf("Current Platform: %v", profile.ProfileChanges[0].Profile.Stats.Attributes.CurrentMtxPlatform).
//				SetTimestamp(time.Now()).
//				Build()).
//			// AddFile("image.png", "Profile operation response.", &b).
//			Build(),
//		)
//
//		return err
//	},
//	LoginRequired:     true,
//	EphemeralResponse: false,
//}
//
//var VbuckPlatformsFriendly = map[string]string{
//	"PSN":         "PlayStation",
//	"Live":        "Xbox",
//	"IOSAppStore": "IOS",
//}
