package commands

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var logout = Command{
	Create: discord.SlashCommandCreate{
		Name:        "logout",
		Description: "Log out of one (or all) of your saved Epic Games accounts.",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionBool{
				Name:        "bulk",
				Description: "Log out of all of your saved Epic Games accounts.",
			},
		},
	},
	Handler: func(event *events.ApplicationCommandInteractionCreate) error {
		return fmt.Errorf("not implemented")
	},
}
