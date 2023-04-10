package commands

import (
	"blast/api"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
)

type CommandHandler func(event *events.ApplicationCommandInteractionCreate) error

type Command struct {
	Create  discord.SlashCommandCreate
	Handler CommandHandler
}

type CommandList []discord.ApplicationCommandCreate

var list = CommandList{}
var handlers = map[string]CommandHandler{}

var blast = api.New()

func (c Command) register() {
	list = append(list, c.Create)
	handlers[c.Create.Name] = c.Handler
}

func Setup() CommandList {
	auth.register()
	login.register()
	logout.register()
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

func FriendlyEmbed(e *discord.EmbedBuilder) *discord.EmbedBuilder {
	return e.SetColor(0xFB5A32).SetTimestamp(time.Now())
}

type MessageDelete struct {
	RawMessage discord.Message
	Timestamp  time.Time
}

var DeletedMessages = make(map[snowflake.ID][]MessageDelete)
