package commands

import (
	"blast/common"
	"blast/db"
	vinderman "github.com/0xDistrust/Vinderman"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var Commands = []discord.ApplicationCommandCreate{
	account,
	//auth,
	//daily,
	//// ephemeral,
	//friends,
	//invite,
	//launch,
	//locker,
	//login,
	//logout,
	//mcp,
	//mnemonic,
	//// offers,
	ping,
	//skipTutorial,
	//vbucks,
	//// test,
	//
	//template,
}

type CommandHandler func(event *handler.CommandEvent, client *common.Client, user db.UserEntry, credentials vinderman.UserCredentials, data discord.SlashCommandInteractionData) error

type Command struct {
	Handler           CommandHandler
	LoginRequired     bool
	EphemeralResponse bool
}
