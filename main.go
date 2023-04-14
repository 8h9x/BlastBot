package main

import (
	"blast/commands"
	"blast/db"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	if err := loadEnv(); err != nil {
		log.Fatal("Unable to load local .env file.")
	}

	if err := db.Init(os.Getenv("MONGO_URI"), "blast"); err != nil {
		log.Fatal("error while connecting to database: ", err)
	}

	client, err := disgo.New(os.Getenv("DISCORD_TOKEN"),
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuilds,
				gateway.IntentGuildMessages,
				gateway.IntentDirectMessages,
			),
		),
		bot.WithEventListenerFunc(func(e *events.Ready) {
			client := e.Client()

			list := commands.Setup()

			if _, err := client.Rest().SetGuildCommands(client.ApplicationID(), 1032441072938405898, list); err != nil {
				log.Fatal("error while registering commands: ", err)
			}

			e.Client().SetPresence(context.Background(), gateway.WithCompetingActivity("UwU"))

			log.Println("Blast is ready!")
		}),
		bot.WithEventListenerFunc(func(event *events.ApplicationCommandInteractionCreate) {
			event.DeferCreateMessage(false)

			data := event.SlashCommandInteractionData()

			user, err := db.Fetch[db.UserEntry]("users", bson.M{"discordId": event.User().ID.String()})
			if err != nil {
				// panic(err)
			}

			// check user logged in
			if len(user.Accounts) == 0 && commands.RequiresLogin(data.CommandName()) {
				event.Client().Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.NewMessageUpdateBuilder().
					SetEmbeds(discord.NewEmbedBuilder().
						SetDescription("You don't have any accounts saved!").
						Build(),
					).
					Build(),
				)

				return
			}

			if handler, ok := commands.Handler(data.CommandName()); ok {
				go func() {
					err := handler(event, user)
					if err != nil {
						log.Println(err)

						event.Client().Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.NewMessageUpdateBuilder().
							SetEmbeds(discord.NewEmbedBuilder().
								SetDescriptionf("```\nAn error occurred while executing the command: %s\n```", err.Error()).
								Build(),
							).
							Build(),
						)
					}
				}()
			}
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	if err = client.OpenGateway(context.TODO()); err != nil {
		log.Fatal(err)
	}

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s
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
