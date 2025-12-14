package interactions

import (
	"log/slog"
	"time"

	"github.com/8h9x/BlastBot/database/internal/interactions/accounts"
	"github.com/8h9x/BlastBot/database/internal/interactions/cloudstorage"
	"github.com/8h9x/BlastBot/database/internal/interactions/login"
	"github.com/8h9x/BlastBot/database/internal/interactions/logout"
	"github.com/8h9x/BlastBot/database/internal/interactions/mcp"
	"github.com/8h9x/BlastBot/database/internal/interactions/mnemonic"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/handler/middleware"
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

var Logger handler.Middleware = func(next handler.Handler) handler.Handler {
	return func(event *handler.InteractionEvent) error {
		event.Client().Logger().InfoContext(event.Ctx, "handling interaction", slog.Any("interaction", event.Interaction), slog.Any("vars", event.Vars))
		return next(event)
	}
}

func CommandHandlerErrorRespond(event *handler.InteractionEvent, err error) {
	slog.Error("interation handling error:", err.Error())

	embed := discord.NewEmbedBuilder().
		SetColor(0xFB2C36).
		SetTimestamp(time.Now()).
		SetTitle("<:exclamation:1096641657396539454> We hit a roadblock!").
		SetDescriptionf("If this issue persists, join our [support server](https://discord.gg/astra-921104988363694130)```\n%s\n```", err.Error())

	event.CreateMessage(discord.NewMessageCreateBuilder().
		SetEmbeds(embed.
			Build(),
		).
		Build(),
	)
}

func init() {
	Router.Use(Logger)
	Router.Use(middleware.GoErr(CommandHandlerErrorRespond))

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
	RegisterCommand(cloudstorage.Definition,
		[]Command{
			{
				Pattern: "/cloudstorage",
				Handler: accounts.AddHandler,
			},
		}...,
	)
	RegisterCommand(login.Definition, Command{
		Pattern: "/login",
		Handler: login.Handler,
	})
	RegisterCommand(login.SessionCheckDefinition, Command{
		Pattern: "/sessioncheck",
		Handler: login.SessionCheckHandler,
	})
	RegisterCommand(logout.Definition, Command{
		Pattern: "/logout",
		Handler: logout.Handler,
	})
	RegisterCommand(mcp.Definition, Command{
		Pattern: "/mcp",
		Handler: mcp.Handler,
	})
	RegisterCommand(mnemonic.Definition,
		[]Command{
			{
				Pattern: "/mnemonic/info",
				Handler: mnemonic.InfoHandler,
			},
		}...,
	)
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
