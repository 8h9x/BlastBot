package commands

import (
	"blast/api"
	"blast/db"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
)

type CommandHandler func(event *events.ApplicationCommandInteractionCreate, user db.UserEntry) error

type Command struct {
	Create        discord.SlashCommandCreate
	Handler       CommandHandler
	LoginRequired bool
}

type CommandList []discord.ApplicationCommandCreate

var list = CommandList{}
var handlers = map[string]CommandHandler{}
var cmds = map[string]Command{}

var blast = api.New()

func (c Command) register() {
	list = append(list, c.Create)
	cmds[c.Create.Name] = c
	handlers[c.Create.Name] = c.Handler
}

func Setup() CommandList {
	account.register()
	auth.register()
	login.register()
	logout.register()
	mcp.register()
	mnemonic.register()

	return list
}

func List() CommandList {
	return list
}

func Handler(name string) (CommandHandler, bool) {
	if handler, ok := handlers[name]; ok {
		return handler, true
	}

	return nil, false
}

func RequiresLogin(name string) bool {
	return cmds[name].LoginRequired
}

func FriendlyEmbed(e *discord.EmbedBuilder) *discord.EmbedBuilder {
	return e.SetColor(0xFB5A32).SetTimestamp(time.Now())
}

type MessageDelete struct {
	RawMessage discord.Message
	Timestamp  time.Time
}

var DeletedMessages = make(map[snowflake.ID][]MessageDelete)
