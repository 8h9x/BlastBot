package common

import (
	"github.com/0xDistrust/Vinderman"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/log"
	"github.com/disgoorg/paginator"
)

type Client struct {
	Discord        *bot.Client
	Epic           *vinderman.Client
	Paginator      *paginator.Manager
	CommandHandler *handler.Mux
	Log            *log.SimpleLogger
}
