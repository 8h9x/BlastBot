package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/log"
	"github.com/disgoorg/paginator"
	"github.com/disgoorg/snowflake/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"

	"blast/api"
	"blast/api/consts"
	"blast/commands"
	"blast/db"
)

func main() {
	if err := loadEnv(); err != nil {
		log.Error("Unable to load local .env file.")
	}

	if err := db.Init(os.Getenv("MONGO_URI"), "blast"); err != nil {
		log.Error("error while connecting to database: ", err)
	}

	logger := log.New(log.Ldate | log.Ltime | log.Lshortfile)

	h := handler.New()

	// h.Command("/test", commands.TestHandler)
	// h.Command("/login", commands.LoginHandler)
	// h.Autocomplete("/test", commands.TestAutocompleteHandler)
	// h.Component("test_button", components.TestComponent)
	h.Command("/account/info", CommandHandlerWrapper(commands.AccountInfo))
	h.Command("/auth/bearer", CommandHandlerWrapper(commands.AuthBearer))
	h.Command("/auth/client", CommandHandlerWrapper(commands.AuthClient))
	h.Command("/auth/device", CommandHandlerWrapper(commands.AuthDevice))
	h.Command("/auth/exchange", CommandHandlerWrapper(commands.AuthExchange))
	h.Command("/login", CommandHandlerWrapper(commands.Login))
	h.Command("/logout", CommandHandlerWrapper(commands.Logout))
	h.Command("/mcp", CommandHandlerWrapper(commands.MCP))
	h.Command("/mnemonic/favorites/add", CommandHandlerWrapper(commands.MnemonicFavoritesAdd))
	h.Command("/mnemonic/favorites/list", CommandHandlerWrapper(commands.MnemonicFavoritesList))
	h.Command("/mnemonic/favorites/remove", CommandHandlerWrapper(commands.MnemonicFavoritesRemove))
	h.Command("/mnemonic/info", CommandHandlerWrapper(commands.MnemonicInfo))

	client, err := disgo.New(os.Getenv("DISCORD_TOKEN"),
		bot.WithLogger(logger),
		bot.WithGatewayConfigOpts(gateway.WithIntents(gateway.IntentGuilds, gateway.IntentGuildMessages)),
		bot.WithCacheConfigOpts(cache.WithCaches(cache.FlagGuilds)),
		bot.WithEventListeners(h, paginator.New()),
		bot.WithEventListenerFunc(func(e *events.Ready) {
			logger.Info("Blast! is ready!")
		}),
	)
	if err != nil {
		logger.Fatal("Failed to setup blast: ", err)
	}

	if _, exists := os.LookupEnv("PROD"); exists {
		logger.Info("Syncing global commands")
		_, err = client.Rest().SetGlobalCommands(client.ApplicationID(), commands.Commands)
		if err != nil {
			logger.Errorf("Failed to sync commands: %s", err)
		}
	} else {
		logger.Infof("Syncing dev (%s) commands", os.Getenv("DISCORD_DEV_GUILD"))
		_, err = client.Rest().SetGuildCommands(client.ApplicationID(), snowflake.GetEnv("DISCORD_DEV_GUILD"), commands.Commands)
		if err != nil {
			logger.Errorf("Failed to sync commands: %s", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = client.OpenGateway(ctx); err != nil {
		logger.Errorf("Failed to connect to gateway: %s", err)
	}
	defer client.Close(context.TODO())

	logger.Info("Blast! is running. Press CTRL-C to exit.")

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s

	logger.Info("Shutting down...")
}

func CommandHandlerWrapper(c commands.Command) handler.CommandHandler {
	return func(event *handler.CommandEvent) error {
		event.DeferCreateMessage(c.EphemeralResponse)

		user, err := db.Fetch[db.UserEntry]("users", bson.M{"discordId": event.User().ID.String()})
		if err != nil {
			event.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().SetContent("Database query error!").Build())
		}

		blast := api.New()

		credentials := api.UserCredentialsResponse{}

		if len(user.Accounts) == 0 && c.LoginRequired {
			CommandHandlerErrorRespond(event, fmt.Errorf("you don't have any accounts saved"))
		} else if c.LoginRequired {
			credentials, err = blast.RefreshTokenLogin(consts.FORTNITE_PC_CLIENT_ID, consts.FORTNITE_PC_CLIENT_SECRET, user.Accounts[user.SelectedAccount].RefreshToken)
			if err != nil {
				CommandHandlerErrorRespond(event, err)
			}
		}

		err = c.Handler(event, *blast, user, credentials, event.SlashCommandInteractionData())

		if err != nil {
			CommandHandlerErrorRespond(event, err)
		}

		return nil
	}
}

func CommandHandlerErrorRespond(event *handler.CommandEvent, err error) {
	event.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
		SetEmbeds(discord.NewEmbedBuilder().
			SetDescriptionf("```\nAn error occurred while executing the command: %s\n```", err.Error()).
			Build(),
		).
		Build(),
	)
}

// only load .env file if not prod
func loadEnv() error {
	_, isProd := os.LookupEnv("PROD")

	if !isProd {
		err := godotenv.Load(".env")
		if err != nil {
			return err
		}
	}

	return nil
}
