package test

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var Definition = discord.SlashCommandCreate{
	Name:        "componenttest",
	Description: "Command to send kitchen sink component message for feature testing",
}

func Handler(event *handler.CommandEvent) error {
	err := event.CreateMessage(discord.MessageCreate{
		Flags: discord.MessageFlagIsComponentsV2,
		Components: []discord.LayoutComponent{
			discord.NewContainer(
				discord.NewTextDisplay("## My Locker Loadout"),
				discord.NewTextDisplay("# CHARACTER"),
				discord.NewActionRow(
					discord.NewStringSelectMenu("select-section", "Preset 1",
						discord.NewStringSelectMenuOption("Character", "section-character"),
						discord.NewStringSelectMenuOption("Emotes", "section-emotes"),
						discord.NewStringSelectMenuOption("Sidekicks", "section-sidekicks"),
						discord.NewStringSelectMenuOption("Wraps", "section-wraps"),
						discord.NewStringSelectMenuOption("Lobby", "section-lobby"),
						discord.NewStringSelectMenuOption("Cars", "section-cars"),
						discord.NewStringSelectMenuOption("Instruments", "section-instruments"),
						discord.NewStringSelectMenuOption("Music", "section-music"),
					),
				),
				discord.NewActionRow(
					discord.NewButton(discord.ButtonStyleSecondary, "Rename", "prev-ssdlkfdsfnection", "", 0),
					discord.NewButton(discord.ButtonStyleDanger, "Delete", "next-dsfogkdsjfnsection", "", 0),
				),
			),
			discord.NewContainer(
				discord.NewActionRow(
					discord.NewStringSelectMenu("select-section2", "Outfit",
						discord.NewStringSelectMenuOption("Character", "section-character"),
						discord.NewStringSelectMenuOption("Emotes", "section-emotes"),
						discord.NewStringSelectMenuOption("Sidekicks", "section-sidekicks"),
						discord.NewStringSelectMenuOption("Wraps", "section-wraps"),
						discord.NewStringSelectMenuOption("Lobby", "section-lobby"),
						discord.NewStringSelectMenuOption("Cars", "section-cars"),
						discord.NewStringSelectMenuOption("Instruments", "section-instruments"),
						discord.NewStringSelectMenuOption("Music", "section-music"),
					),
				),
				discord.NewSmallSeparator(),
				discord.SectionComponent{
					Components: []discord.SectionSubComponent{
						discord.NewTextDisplay("### OUTFIT"),
						discord.NewTextDisplay("Iconic Kim Kardashian"),
					},
					Accessory: discord.NewThumbnail("https://fortnite-api.com/images/cosmetics/br/Character_QuicheLorraineLime/icon.png"),
				},
				discord.NewActionRow(
					discord.NewButton(discord.ButtonStyleSuccess, "Change", "next-sectidon", "", 0),
					discord.NewButton(discord.ButtonStylePrimary, "Random", "next-sectdidon", "", 0),
				),
			),
		},
	})

	return err
}
