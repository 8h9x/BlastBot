package main

import (
	"context"
//	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/8h9x/BlastBot/internal/database"
	"github.com/8h9x/BlastBot/internal/interactions"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	_, isProd := os.LookupEnv("PROD")
	if !isProd {
		err := godotenv.Load(".env")
		if err != nil {
			slog.Error(fmt.Sprintf("error loading .env file: %s", err))
			os.Exit(1)
		}
	}


//	url := fmt.Sprintf("%s?authToken=%s", os.Getenv("TURSO_DATABASE_URL"), os.Getenv("TURSO_AUTH_TOKEN"))

//	db, err := sql.Open("libsql", url)
//	if err != nil {
//		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", url, err)
//		os.Exit(1)
//	}
//  	defer db.Close()
//
//	log.Println(db.Stats().MaxOpenConnections)

	if err := database.Init(os.Getenv("MONGODB_URI"), "blast"); err != nil {
		slog.Error(fmt.Sprintf("error while connecting to database: %s", err))
	}

	client, err := disgo.New(os.Getenv("DISCORD_BOT_TOKEN"),
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuilds,
				gateway.IntentDirectMessages,
			),
			gateway.WithPresenceOpts(
				gateway.WithOnlineStatus(discord.OnlineStatusOnline),
				gateway.WithCustomActivity("Being the best utility bot"),
			),
		),
		bot.WithEventListenerFunc(func(e *events.Ready) {
			slog.Info("Connected to discord gateway")
		}),
	)
	if err != nil {
		panic(err)
	}

	if err = client.OpenGateway(context.TODO()); err != nil {
		panic(err)
	}

	guildID := snowflake.ID(0)
	if _, exists := os.LookupEnv("PROD"); !exists {
		guildID = snowflake.GetEnv("DISCORD_DEV_GUILD")
		if guildID == 0 {
			log.Fatal("unable to sync commands because environment PROD variable is not present and DISCORD_DEV_GUILD enviroment variable is not set")
		}
	}

	commandEnvString := "global"
	if guildID > 0 {
		commandEnvString = guildID.String()
	}

	slog.Info(fmt.Sprintf("Syncing (%s) commands", commandEnvString))

	err = interactions.SyncCommands(client, guildID)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to sync (%s) commands: %s", commandEnvString, err))
	}

	slog.Info(fmt.Sprintf("Finished syncing (%s) commands", commandEnvString))

	client.AddEventListeners(interactions.Router)

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s
}
