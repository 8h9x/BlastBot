package commands

import (
	"blast/api"
	"blast/api/consts"
	"blast/db"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var account = Command{
	Create: discord.SlashCommandCreate{
		Name:        "account",
		Description: "Account related commands.",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionSubCommand{
				Name:        "info",
				Description: "Retrieves general account information.",
			},
		},
	},
	Handler: func(event *events.ApplicationCommandInteractionCreate, user db.UserEntry) error {
		data := event.SlashCommandInteractionData()

		refreshCredentials, err := blast.RefreshTokenLogin(consts.FORTNITE_PC_CLIENT_ID, consts.FORTNITE_PC_CLIENT_SECRET, user.Accounts[user.SelectedAccount].RefreshToken)
		if err != nil {
			return err
		}

		account, err := blast.FetchAccountInformation(refreshCredentials)
		if err != nil {
			return err
		}

		brInfo, err := blast.FetchAccountBRInfo(refreshCredentials)
		if err != nil {
			return err
		}

		athenaProfile, err := blast.ProfileOperationStr(refreshCredentials, "QueryProfile", "athena", "{}")
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

		// XP: response.data.profileChanges[0].profile.stats.attributes.xp

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

		switch *data.SubCommandName {
		case "info":
			embed := discord.NewEmbedBuilder().
				SetAuthorIconf("https://fortnite-api.com/images/cosmetics/br/%s/icon.png", strings.Replace(locker.Attributes.LockerSlotsData.Slots["Character"].Items[0], "AthenaCharacter:", "", -1)). // TODO set author icon to bot user avatar
				SetColor(0xFB5A32).
				SetTimestamp(time.Now()).
				SetAuthorNamef("%s | %s", account.DisplayName, refreshCredentials.AccountId).
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
				// SetDescription(refreshCredentials.AccountId).
				Build()

			_, err = event.Client().Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.NewMessageUpdateBuilder().SetEmbeds(embed).Build())
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown subcommand")
		}

		return nil
	},
	LoginRequired: true,
}

func boolToEmoji(b bool) string {
	if b {
		return "<a:thumbs_up:1096490549218906193>"
	}

	return "<a:thumbs_down:1096490737715130388>"
}

// ApplicationCommandOptionSubCommand
