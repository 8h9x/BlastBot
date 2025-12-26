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
	err := event.CreateMessage(discord.MessageCreate{
		Flags: discord.MessageFlagIsComponentsV2,
		Components: []discord.LayoutComponent{
			discord.NewContainer(
				discord.SectionComponent{
					Components: []discord.SectionSubComponent{
						discord.NewTextDisplay("### OUTFIT"),
						discord.NewTextDisplay("Iconic Kim Kardashian"),
					},
					Accessory: discord.NewThumbnail("https://fortnite-api.com/images/cosmetics/br/Character_QuicheLorraineLime/icon.png"),
				},
				discord.NewActionRow(
					discord.NewButton(discord.ButtonStyleSuccess, "Equip", "next-sectidon", "", 0),
					discord.NewButton(discord.ButtonStylePrimary, "Favorite", "next-sectdidon", "", 0),
					discord.NewButton(discord.ButtonStylePrimary, "Purchase", "next-secteedidon", "", 0),
					discord.NewButton(discord.ButtonStylePrimary, "Wishlist", "next-seeectdidon", "", 0),
					),
			),
		},
	})

	return err
}
