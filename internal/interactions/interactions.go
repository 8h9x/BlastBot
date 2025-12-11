package interactions

import (
	"github.com/8h9x/BlastBot/database/internal/interactions/accounts"
	"github.com/8h9x/BlastBot/database/internal/interactions/login"
	"github.com/8h9x/BlastBot/database/internal/interactions/logout"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
)

var (
	Router      = handler.New()
	definitions []discord.ApplicationCommandCreate
)

type Command struct {
	Pattern string
	Handler handler.CommandHandler
}

func init() {
	RegisterCommand(accounts.Definition,
		[]Command{
			{
				Pattern: "/accounts/add",
				Handler: accounts.AddHandler,
			},
			{
				Pattern: "/accounts/remove",
				Handler: accounts.RemoveHandler,
			},
			{
				Pattern: "/accounts/switch",
				Handler: accounts.SwitchHandler,
			},
		}...,
	)
	RegisterCommand(login.Definition, Command{
		Pattern: "/login",
		Handler: login.Handler,
	})
	RegisterCommand(logout.Definition, Command{
		Pattern: "/logout",
		Handler: logout.Handler,
	})
}

func RegisterCommand(def discord.ApplicationCommandCreate, cmds ...Command) {
	definitions = append(definitions, def)
	for _, cmd := range cmds {
		Router.Command(cmd.Pattern, cmd.Handler)
	}
}

func SyncCommands(client bot.Client, guildID snowflake.ID) error {
	if guildID.String() == "0" {
		_, err := client.Rest().SetGlobalCommands(client.ApplicationID(), definitions)
		if err != nil {
			return err
		}
	} else {
		_, err := client.Rest().SetGuildCommands(client.ApplicationID(), guildID, definitions)
		if err != nil {
			return err
		}
	}

	return nil
}
