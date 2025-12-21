package winterfest

import (
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

type rewardGraph struct {
	TemplateID  string
	UnlockEpoch int64
	Nodes       map[string]rewardNode `json:"nodes"`
}

type rewardNode struct {
	TagName string `json:"tagName"`
	Items   []string `json:"items"`
}

var RewardGraph = rewardGraph{
	TemplateID: "AthenaRewardGraph:s39_winterfest",
	UnlockEpoch: time.Date(2025, 12, 17, 14, 0, 0, 0, time.UTC).UnixMilli(),
	Nodes: map[string]rewardNode{
		"ERG.Node.A.1": {
			Items: []string{"Character_ChubbyJingle"},
		},
		"ERG.Node.A.2": {
			Items: []string{"Character_FrostIron"},
		},
		"ERG.Node.A.3": {
			Items: []string{"Backpack_PolarGander"},
		},
		"ERG.Node.A.4": {
			Items: []string{"Emoticon_S39Winterfest", "Spray_WinterFest3"},
		},
		"ERG.Node.A.5": {
			Items: []string{"Backpack_WinkFin"},
		},
		"ERG.Node.A.6": {
			Items: []string{"SID_Placeholder_728"},
		},
		"ERG.Node.A.7": {
			Items: []string{"Pickaxe_HunSting"},
		},
		"ERG.Node.A.8": {
			Items: []string{"Wrap_ShavedIce"},
		},
		"ERG.Node.A.9": {
			Items: []string{"Contrail_HollowStun"},
		},
		"ERG.Node.A.10": {
			Items: []string{"Contrail_HollowStun"},
		},
		"ERG.Node.A.11": {
			Items: []string{"Sparks_RainbowHouse_Guitar"},
		},
		"ERG.Node.A.12": {
			Items: []string{"Sparks_FrostIron_Bass"},
		},
		"ERG.Node.A.13": {
			Items: []string{"Glider_RumpleWisp"},
		},
		"ERG.Node.B.1": {
			Items: []string{"Spray_WinterFest", "Spray_WinterFest2"},
		},
	},
}

var Definition = discord.SlashCommandCreate{
	Name:        "winterfest",
	Description: "Command to send kitchen sink component message for feature testing",
}

func Setup(router *handler.Mux) {
	router.Command("/winterfest", func(event *handler.CommandEvent) error {
		return nil
	})

	router.Component("/winterfest-cabin-prev", func(e *handler.ComponentEvent) error {
		return nil
	})

	router.Component("/winterfest-cabin-next", func(e *handler.ComponentEvent) error {
		return nil
	})
}

func Handler(event *handler.CommandEvent) error {
	err := event.CreateMessage(discord.MessageCreate{
		Flags: discord.MessageFlagIsComponentsV2,
		Components: []discord.LayoutComponent{
			discord.NewContainer(
				discord.NewTextDisplay("## <:snowflake:1451209367738388532> Chapter 7, Season 1 Winterest Cabin"),
				discord.NewTextDisplayf("Happy holidays, **%s**. You have **%s** unopened presents.", "8h9x", "999"),
				discord.NewLargeSeparator(),
				discord.SectionComponent{
					Components: []discord.SectionSubComponent{
						discord.NewTextDisplay("### Cheeks"),
						discord.NewTextDisplay("OUTFIT\n`ERG.Node.A.1`"),
					},
					Accessory: discord.NewThumbnail("https://fortnite-api.com/images/cosmetics/br/Character_ChubbyJingle/icon.png"),
				},
				discord.NewActionRow(
					discord.NewButton(discord.ButtonStyleSuccess, "Open", "winterfest-claim-character_chubbyjingle", "", 0).
						WithEmoji(discord.NewCustomComponentEmoji(1451205113669484655)),
				),
				discord.NewSmallSeparator(),
				discord.SectionComponent{
					Components: []discord.SectionSubComponent{
						discord.NewTextDisplay("### Glacial Dummy"),
						discord.NewTextDisplay("OUTFIT\n`ERG.Node.A.2`"),
					},
					Accessory: discord.NewThumbnail("https://fortnite-api.com/images/cosmetics/br/Character_FrostIron/icon.png"),
				},
				discord.NewActionRow(
					discord.NewButton(discord.ButtonStyleSuccess, "Open", "winterfest-claim-character_frostiron", "", 0).
						WithEmoji(discord.NewCustomComponentEmoji(1451205113669484655)),
				),
				discord.NewLargeSeparator(),
				discord.NewActionRow(
					discord.NewPrimaryButton("", "/winterfest-cabin-prev").
						WithEmoji(discord.NewCustomComponentEmoji(1451208271984722073)),
					discord.NewSecondaryButton("1 / 7", "/winterfest-cabin-page-num"),
					discord.NewPrimaryButton("", "/winterfest-cabin-next").
						WithEmoji(discord.NewCustomComponentEmoji(1451208303370702941)),
				),
			),
		},
	})

	return err
}
