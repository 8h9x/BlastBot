package commands

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var logout = Command{
	Create: discord.SlashCommandCreate{
		Name:        "logout",
		Description: "wip",
	},
	Handler: func(event *events.ApplicationCommandInteractionCreate) error {
		return fmt.Errorf("not implemented")
	},
}
