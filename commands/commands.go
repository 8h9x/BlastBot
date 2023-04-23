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
	daily,
	// ephemeral,
	friends,
	invite,
	launch,
	locker,
	login,
	logout,
	mcp,
	mnemonic,
	// offers,
	ping,
	skipTutorial,
	vbucks,
	// test,
}

type CommandHandler func(event *handler.CommandEvent, blast api.EpicClient, user db.UserEntry, credentials api.UserCredentialsResponse, data discord.SlashCommandInteractionData) error

type Command struct {
	Handler           CommandHandler
	LoginRequired     bool
	EphemeralResponse bool
}
