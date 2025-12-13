package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/8h9x/BlastBot/database/internal/database"
	"github.com/8h9x/BlastBot/database/internal/interactions"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/bot"
	"github.com/joho/godotenv"
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

	if err := database.Init(os.Getenv("MONGODB_URI"), "blast"); err != nil {
		slog.Error(fmt.Sprintf("error while connecting to database: %s", err))
	}

	commandHandler := handler.New()
	commandHandler.Command("/login", dummyCommand)                     // Alias for /accounts/add method: DeviceCode client: FORTNITE_PS4_US_CLIENT
	commandHandler.Command("/logout", dummyCommand)                    // Alias for /accounts/remove {currentAccount} -- also need bulk logout
	commandHandler.Command("/compose_mcp_request", dummyCommand)       // Create your own raw MCP request with guided components
	commandHandler.Command("/launch", dummyCommand)                    // Generate launch args
	commandHandler.Command("/accounts/add", dummyCommand)              // Add an account using any method and client
	commandHandler.Command("/accounts/status", dummyCommand)           // Displays brief information of all accounts, number of accounts, emphasizes currently selected
	commandHandler.Command("/accounts/switch", dummyCommand)           // Swap between synced accounts
	commandHandler.Command("/accounts/remove", dummyCommand)           // Remove an account from management and delete it's database entry
	commandHandler.Command("/friends/add", dummyCommand)               // Add a friend :)
	commandHandler.Command("/friends/list", dummyCommand)              // List of current friends
	commandHandler.Command("/friends/remove", dummyCommand)            // Remove a friend :(
	commandHandler.Command("/friends/requests/list", dummyCommand)     // List of current incoming friend requests
	commandHandler.Command("/friends/requests/accept", dummyCommand)   // Accept an incoming friend request
	commandHandler.Command("/friends/requests/decline", dummyCommand)  // Decline an incoming friend request
	commandHandler.Command("/playlist/favorites/add", dummyCommand)    // Favorite a playlist (mnemonic)
	commandHandler.Command("/playlist/favorites/list", dummyCommand)   // List of all favorite playlists (mnemonics)
	commandHandler.Command("/playlist/favorites/remove", dummyCommand) // Unfavorite a playlist (mnemonic)
	commandHandler.Command("/playlist/recents", dummyCommand)          // List of all recently played playlists
	commandHandler.Command("/playlist/info", dummyCommand)             // Fetch playlist meta information
	commandHandler.Command("/locker/image", dummyCommand)              // Generate an image of locker data
	commandHandler.Command("/locker/equip", dummyCommand)              // Equip a cosmetic item
	commandHandler.Command("/locker/loadouts/select", dummyCommand)    // Change active loadout
	commandHandler.Command("/locker/loadouts/list", dummyCommand)      // List loadouts
	commandHandler.Command("/lobby/equip", dummyCommand)               // Temporarily equip ANY cosmetic item in the lobby (only can be seen by peers, only works in the lobby--duh!)
	commandHandler.Command("/lobby/crowns", dummyCommand)              // Temporarily display an arbitrary number of crowns in the lobby (only can be seen by peers, only works in the lobby--duh!)
	commandHandler.Command("/party/invite", dummyCommand)              // Sends a party invite
	commandHandler.Command("/party/kick", dummyCommand)                // Kick someone from your party
	commandHandler.Command("/party/leave", dummyCommand)               // Leave your current party
	commandHandler.Command("/profile", dummyCommand)                   // Return file from QueryProfile data for the inputted profile_id

	commandHandler.Command("/auto/research", dummyCommand) // Use research points in stw
	commandHandler.Command("/cloudstorage", dummyCommand)  // For downloading/uploading/editing game settings files

	client, err := disgo.New(os.Getenv("DISCORD_BOT_TOKEN"),
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuilds,
				gateway.IntentDirectMessages,
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

	client.AddEventListeners(interactions.Router)

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s
}

func dummyCommand(e *handler.CommandEvent) error {
	return nil
}
