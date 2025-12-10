package commands

import (
	"github.com/8h9x/BlastBot/database/internal/database"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var Commands = []discord.ApplicationCommandCreate{
	login,
}

type CommandHandler func(event *handler.CommandEvent, user database.UserEntry, data discord.SlashCommandInteractionData) error

type Commander interface {
	Handler() error
}

type Command struct {
	Handler           CommandHandler
	LoginRequired     bool
	EphemeralResponse bool
}
