package commands

import (
	"blast/api"
	"blast/db"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var Commands = []discord.ApplicationCommandCreate{
	account,
	auth,
	login,
	logout,
	mcp,
	mnemonic,
	// test,
}

type CommandHandler func(event *handler.CommandEvent, blast api.EpicClient, user db.UserEntry, credentials api.UserCredentialsResponse, data discord.SlashCommandInteractionData) error

type Command struct {
	Handler           CommandHandler
	LoginRequired     bool
	EphemeralResponse bool
}
