package commands

import (
	"blast/api"
	"blast/api/consts"
	"blast/db"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"go.mongodb.org/mongo-driver/bson"
)

var account = discord.SlashCommandCreate{
	Name:        "account",
	Description: "Account related commands.",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionSubCommand{
			Name:        "info",
			Description: "Retrieves general account information.",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "switch",
			Description: "Switches your selected epic games account.",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "vbucks",
			Description: "Check how many vbucks are on your account(s).",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionBool{
					Name:        "bulk",
					Description: "Whether to check all accounts or just the selected one.",
				},
			},
		},
	},
}

var AccountInfo = Command{
	Handler: func(event *handler.CommandEvent, blast api.EpicClient, user db.UserEntry, credentials api.UserCredentialsResponse, data discord.SlashCommandInteractionData) error {
		account, err := blast.FetchMyAccountInfo(credentials)
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

		avatarURL, err := blast.FetchAvatarURL(credentials)
		if err != nil {
			return err
		}

		embed := discord.NewEmbedBuilder().
			SetAuthorIcon(avatarURL).
			SetColor(0xFB5A32).
			SetTimestamp(time.Now()).
			SetAuthorNamef("%s | %s", account.DisplayName, credentials.AccountID).
			AddField("<:llama:1096476378121126000> Account Level", fmt.Sprint(attributes.AccountLevel), true).
			AddField(fmt.Sprintf("%s Season Level", bpEmoji), fmt.Sprint(attributes.Level), true).
			AddField("<:battle_star:1096473613504368640> Battlestars", fmt.Sprint(attributes.Battlestars), true).
			AddField("<:xp:1096486771820347453> Season XP", fmt.Sprint(attributes.XP), true).
			AddField("<:goldbars:1096469831034863637> Bars", fmt.Sprint(brInfo.Stash.Globalcash), true).
			AddField("<:victory:1096481674872770591> Lifetime Wins", fmt.Sprint(attributes.LifetimeWins), true).
			// AddField("<:victory_crown:1096481681575260314> Crown Wins", "coming soon", true).
			AddField("<:outfit:1096486172655636561> Skin Count", fmt.Sprint(skinCount), true).
			AddField("<:lock:1096491497987260578> MFA Reward Claimed", BoolToEmoji(attributes.MFARewardClaimed), true).
			AddField(fmt.Sprintf(":flag_%s: Region", strings.ToLower(account.Country)), account.Country, true).
			Build()

		_, err = event.Client().Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.NewMessageUpdateBuilder().SetEmbeds(embed).Build())
		return err
	},
	LoginRequired:     true,
	EphemeralResponse: false,
}

var AccountSwitch = Command{
	Handler: func(event *handler.CommandEvent, blast api.EpicClient, user db.UserEntry, credentials api.UserCredentialsResponse, data discord.SlashCommandInteractionData) error {
		switchOptions := make([]discord.StringSelectMenuOption, 0)

		for i, account := range user.Accounts {
			refreshCredentials, err := blast.RefreshTokenLogin(consts.FORTNITE_PC_CLIENT_ID, consts.FORTNITE_PC_CLIENT_SECRET, account.RefreshToken)
			if err != nil {
				if err.(*api.RequestError).Raw.ErrorCode == "errors.com.epicgames.account.auth_token.invalid_refresh_token" {
					col := db.GetCollection("users")
					_, err = col.UpdateOne(context.Background(), bson.M{"discordId": event.User().ID}, bson.M{"$pull": bson.M{"accounts": bson.M{"accountId": account.AccountID}}})
					if err != nil {
						return err
					}

					continue
				}
				return err
			}

			accountInfo, err := blast.FetchMyAccountInfo(refreshCredentials)
			if err != nil {
				return err
			}

			switchOptions = append(switchOptions, discord.NewStringSelectMenuOption(accountInfo.DisplayName, fmt.Sprint(i)).WithDefault(i == user.SelectedAccount))
		}

		_, err := event.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
			AddActionRow(discord.NewStringSelectMenu("switch_account_select", "Select the saved account to switch to", switchOptions...)).
			AddActionRow(discord.NewDangerButton("Cancel", "cancel").WithEmoji(discord.ComponentEmoji{
				Name: "x_",
				ID:   1096630553689739385,
			})).
			Build(),
		)
		return err
	},
	LoginRequired:     true,
	EphemeralResponse: false,
}

var AccountVbucks = Command{
	Handler: func(event *handler.CommandEvent, blast api.EpicClient, user db.UserEntry, credentials api.UserCredentialsResponse, data discord.SlashCommandInteractionData) error {
		if data.Bool("bulk") {
			embed := discord.NewEmbedBuilder().SetColor(0xFB5A32).SetTimestamp(time.Now()).SetTitle("Saved Accounts V-Bucks")

			bulkTotal := 0

			for _, account := range user.Accounts {
				refreshCredentials, err := blast.RefreshTokenLogin(consts.FORTNITE_PC_CLIENT_ID, consts.FORTNITE_PC_CLIENT_SECRET, account.RefreshToken)
				if err != nil {
					if err.(*api.RequestError).Raw.ErrorCode == "errors.com.epicgames.account.auth_token.invalid_refresh_token" {
						col := db.GetCollection("users")
						_, err = col.UpdateOne(context.Background(), bson.M{"discordId": event.User().ID}, bson.M{"$pull": bson.M{"accounts": bson.M{"accountId": account.AccountID}}})
						if err != nil {
							return err
						}

						continue
					}
					return err
				}

				accountInfo, err := blast.FetchMyAccountInfo(refreshCredentials)
				if err != nil {
					return err
				}

				commonCoreProfile, err := blast.ProfileOperationStr(refreshCredentials, "QueryProfile", "common_core", "{}")
				if err != nil {
					continue
					// return err
				}

				defer commonCoreProfile.Close()

				var profile api.CommonCoreProfile
				err = json.NewDecoder(commonCoreProfile).Decode(&profile)
				if err != nil {
					return err
				}

				total := 0
				list := ""

				for _, item := range profile.ProfileChanges[0].Profile.Items {
					if strings.HasPrefix(item.TemplateID, "Currency:Mtx") {
						total += item.Quantity
						bulkTotal += item.Quantity

						if len(user.Accounts) < 10 { // only show the extended list if there are less than 10 accounts
							platform := VbuckPlatformsFriendly[item.Attributes.Platform]
							if platform == "" {
								platform = item.Attributes.Platform
							}

							list += fmt.Sprintf("<:vbuck:1099111673966628905> %d - *%s (%s)*\n", item.Quantity, platform, strings.Replace(item.TemplateID, "Currency:Mtx", "", 1))
						}
					}
				}

				if len(user.Accounts) < 10 {
					embed.AddField(accountInfo.DisplayName, fmt.Sprintf("<:vbuck:1099111673966628905> **%d - Total** \n%s", total, list), true)
				} else {
					embed.AddField(accountInfo.DisplayName, fmt.Sprintf("<:vbuck:1099111673966628905> %d", total), true)
				}
			}

			embed.SetFooterTextf("Total: %d V-Bucks", bulkTotal)

			_, err := event.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().SetEmbeds(embed.Build()).Build())
			return err
		}

		account, err := blast.FetchMyAccountInfo(credentials)
		if err != nil {
			return err
		}

		avatarURL, err := blast.FetchAvatarURL(credentials)
		if err != nil {
			return err
		}

		commonCoreProfile, err := blast.ProfileOperationStr(credentials, "QueryProfile", "common_core", "{}")
		if err != nil {
			return err
		}

		defer commonCoreProfile.Close()

		var profile api.CommonCoreProfile
		err = json.NewDecoder(commonCoreProfile).Decode(&profile)
		if err != nil {
			return err
		}

		total := 0
		list := ""

		for _, item := range profile.ProfileChanges[0].Profile.Items {
			if strings.HasPrefix(item.TemplateID, "Currency:Mtx") {
				total += item.Quantity
				platform := VbuckPlatformsFriendly[item.Attributes.Platform]
				if platform == "" {
					platform = item.Attributes.Platform
				}

				list += fmt.Sprintf("<:vbuck:1099111673966628905> **%d** - %s (%s)\n", item.Quantity, platform, strings.Replace(item.TemplateID, "Currency:Mtx", "", 1))
			}
		}

		_, err = event.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
			SetEmbeds(discord.NewEmbedBuilder().
				SetTitle(fmt.Sprintf("<:vbuck:1099111673966628905> %d Total", total)).
				SetColor(0xFB5A32).
				SetAuthorNamef("%s's V-Bucks", account.DisplayName).
				SetAuthorIcon(avatarURL).
				SetDescription(list).
				SetFooterTextf("Current Platform: %v", profile.ProfileChanges[0].Profile.Stats.Attributes.CurrentMtxPlatform).
				SetTimestamp(time.Now()).
				Build()).
			// AddFile("image.png", "Profile operation response.", &b).
			Build(),
		)

		return err
	},
	LoginRequired:     true,
	EphemeralResponse: false,
}

func BoolToEmoji(b bool) string {
	if b {
		return "<a:thumbs_up:1096490549218906193>"
	}

	return "<a:thumbs_down:1096490737715130388>"
}

var VbuckPlatformsFriendly = map[string]string{
	"PSN":         "PlayStation",
	"Live":        "Xbox",
	"IOSAppStore": "IOS",
}
