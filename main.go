package main

import (
	"blast/commands"
	"blast/common"
	"blast/db"
	"context"
	"errors"
	"github.com/0xDistrust/Vinderman"
	"github.com/0xDistrust/Vinderman/consts"
	"github.com/0xDistrust/Vinderman/request"
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
	"go.mongodb.org/mongo-driver/bson"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := common.LoadEnv(); err != nil {
		log.Error("Unable to load local .env file.")
	}

	if err := db.Init(os.Getenv("MONGO_URI"), "blast"); err != nil {
		log.Error("error while connecting to database: ", err)
	}

	logger := log.New(log.Ldate | log.Ltime | log.Lshortfile)
	h := handler.New()
	manager := paginator.New() // todo: re-add config
	epic := vinderman.New()

	client, err := disgo.New(os.Getenv("DISCORD_TOKEN"),
		bot.WithLogger(logger),
		bot.WithGatewayConfigOpts(gateway.WithIntents(gateway.IntentGuilds, gateway.IntentGuildMessages)),
		bot.WithCacheConfigOpts(cache.WithCaches(cache.FlagGuilds)),
		bot.WithEventListeners(h, manager),
		bot.WithEventListenerFunc(func(e *events.Ready) {
			logger.Info("Connected to discord gateway.")
		}),
	)
	if err != nil {
		logger.Fatal("Failed to setup discord client: ", err)
	}

	blast := &common.Client{
		Discord:        &client,
		Epic:           epic,
		Paginator:      manager,
		CommandHandler: h,
		Log:            logger,
	}

	h.Command("/account/info", CommandHandlerWrapper(blast, commands.AccountInfo, true))
	h.Command("/account/switch", CommandHandlerWrapper(blast, commands.AccountSwitch, true))
	h.Command("/account/vbucks", CommandHandlerWrapper(blast, commands.AccountVbucks, true))
	h.Command("/ping", CommandHandlerWrapper(blast, commands.Ping, true))

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

func CommandHandlerWrapper(client *common.Client, cmd commands.Command, acknowledge bool) handler.CommandHandler {
	return func(event *handler.CommandEvent) error {
		if acknowledge {
			event.DeferCreateMessage(cmd.EphemeralResponse)
		}

		user, err := db.Fetch[db.UserEntry]("users", bson.M{"discordId": event.User().ID})
		if err != nil {
			// event.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().SetContent("Database query error!").Build())
		}

		credentials := vinderman.UserCredentials{}

		if len(user.Accounts) == 0 && cmd.LoginRequired {
			embed := discord.NewEmbedBuilder().
				SetColor(0xFB5A32).
				SetTimestamp(time.Now()).
				SetTitle("<:llama:1096476378121126000> Not the llama you're looking for!").
				SetDescription("You do not have any saved accounts.\nAdd one using the `/login` command.").
				Build()

			_, err := event.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
				SetEmbeds(embed).
				ClearContent().
				Build(),
			)
			return err
		} else if cmd.LoginRequired {
			credentials, err = client.Epic.RefreshTokenLogin(consts.FORTNITE_PC_CLIENT_ID, consts.FORTNITE_PC_CLIENT_SECRET, user.Accounts[user.SelectedAccount].RefreshToken)
			if err != nil {
				if err.(*request.Error[vinderman.EpicErrorResponse]).Raw.ErrorCode == "errors.com.epicgames.account.auth_token.invalid_refresh_token" {
					col := db.GetCollection("users")
					_, err = col.UpdateOne(context.Background(), bson.M{"discordId": event.User().ID}, bson.M{"$pull": bson.M{"accounts": bson.M{"accountId": user.Accounts[user.SelectedAccount].AccountID}}})
					if err != nil {
						return err
					}

					_, err = col.UpdateOne(context.Background(), bson.M{"discordId": user.ID}, bson.M{"$inc": bson.M{"selectedAccount": -1}})
					if err != nil {
						return err
					}

					CommandHandlerErrorRespond(event, errors.New("your account has been removed from the database due to an invalid refresh token. please /login again to add it back"))
				}
				CommandHandlerErrorRespond(event, err)
			}
		}

		go func() {
			if err := cmd.Handler(event, client, user, credentials, event.SlashCommandInteractionData()); err != nil {
				CommandHandlerErrorRespond(event, err)
			}
		}()

		return nil
	}
}

func CommandHandlerErrorRespond(event *handler.CommandEvent, err error) {
	embed := discord.NewEmbedBuilder().
		SetColor(0xFB5A32).
		SetTimestamp(time.Now()).
		SetTitle("<:exclamation:1096641657396539454> We hit a roadblock!").
		SetDescriptionf("If this issue persists, join our [support server](https://discord.gg/astra-921104988363694130)```\n%s\n```", err.Error())

	switch err := err.(type) {
	case *request.Error[vinderman.EpicErrorResponse]:
		embed.SetFooterText(err.Raw.ErrorCode)
	}

	event.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
		SetEmbeds(embed.
			Build(),
		).
		Build(),
	)
}
