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
			data := event.SlashCommandInteractionData()

			if handler, ok := commands.Handler(data.CommandName()); ok {
				go func() {
					err := handler(event)
					if err != nil {
						log.Println(err)

						event.CreateMessage(discord.NewMessageCreateBuilder().
							SetEmbeds(commands.FriendlyEmbed(discord.NewEmbedBuilder().
								SetDescriptionf("```\nAn error occurred while executing the command: %s\n```", err.Error()),
							).Build()).
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
