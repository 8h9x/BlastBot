package cosmetic

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var Definition = discord.SlashCommandCreate{
	Name:        "cosmetic",
	Description: "Command to send kitchen sink component message for feature testing",
}

func Handler(event *handler.CommandEvent) error {
	inline := true

	err := event.CreateMessage(discord.MessageCreate{
		Embeds: []discord.Embed{{
			Title:       "Iconic Kim Kardashian",
			Description: "Mogul. Mom. Icon.",
			Thumbnail: &discord.EmbedResource{
				URL: "https://fortnite-api.com/images/cosmetics/br/Character_QuicheLorraineLime/icon.png",
			},
			Color: 0x5cf2f3,
			Fields: []discord.EmbedField{
				{
					Name:   "ID",
					Value:  "Character_QuicheLorraineLime",
					Inline: &inline,
				},
				{
					Name:   "Type",
					Value:  "Outfit",
					Inline: &inline,
				},
				{
					Name:   "Series",
					Value:  "Icon Series",
					Inline: &inline,
				},
				{
					Name:   "Introduced",
					Value:  "Chapter 7, Season 1.",
					Inline: &inline,
				},
			},
		}},
		Components: []discord.LayoutComponent{
			discord.NewActionRow(
				discord.NewButton(discord.ButtonStyleSuccess, "Equip", "next-sectidon", "", 0),
				discord.NewButton(discord.ButtonStylePrimary, "Favorite", "next-sectdidon", "", 0),
				discord.NewButton(discord.ButtonStylePrimary, "Purchase", "next-secteedidon", "", 0),
				discord.NewButton(discord.ButtonStylePrimary, "Wishlist", "next-seeectdidon", "", 0),
			),
		},
	})

	return err
}
