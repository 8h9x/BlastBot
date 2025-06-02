package main

import (
	"context"
	"fmt"
	"github.com/disgoorg/disgo/handler"
	"github.com/surrealdb/surrealdb.go"
	"github.com/surrealdb/surrealdb.go/pkg/models"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/joho/godotenv"
	"gitlab.com/8h9x/BlastBot/database"
)

func main() {
	_, isProd := os.LookupEnv("PROD")
	if !isProd {
		err := godotenv.Load(".env")
		if err != nil {
			slog.Error(fmt.Sprintf("Error loading .env file: %s", err))
			os.Exit(1)
		}
	}

	// Connect to SurrealDB
	db, err := surrealdb.New(os.Getenv("SURREALDB_CONNECTION_URL"))
	if err != nil {
		panic(err)
	}

	// Set the namespace and database
	if err = db.Use("BlastBot", "integ"); err != nil {
		panic(err)
	}

	// Sign in to authentication `db`
	authData := &surrealdb.Auth{
		Username: os.Getenv("SURREALDB_ADMIN_USERNAME"), // use your setup username
		Password: os.Getenv("SURREALDB_ADMIN_PASSWORD"), // use your setup password
	}
	token, err := db.SignIn(authData)
	if err != nil {
		panic(err)
	}

	// And we can later on invalidate the token if desired
	defer func(token string) {
		if err := db.Invalidate(); err != nil {
			panic(err)
		}
	}(token)

	user1, err := surrealdb.Create[interface{}](db, models.Table("users"), database.User{
		CreatedAt: time.Now(),
		DiscordId: "908900960791834674",
		//EpicAccounts          []models.RecordID `json:"epic_accounts"`
		GlobalFlags:           0,
		SelectedEpicAccountId: "17dcac15e1554c9eb79445a96c859c81",
		UpdatedAt:             time.Now(),
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created user with a struct: %+v\n", user1)

	commandHandler := handler.New()

	commandHandler.Command("/login", dummyCommand)
	commandHandler.Command("/compose_mcp_request", dummyCommand)       // Create your own raw MCP request with guided components
	commandHandler.Command("/launch", dummyCommand)                    // Generate launch args
	commandHandler.Command("/accounts/add", dummyCommand)              // Alias for /login method: DeviceCode client: FORTNITE_PS4_US_CLIENT
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
	commandHandler.Command("/profile", dummyCommand)                   // Return file from QueryProfile data for the inputted profile_id

	client, err := disgo.New(os.Getenv("DISCORD_BOT_TOKEN"),
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuilds,
				// 				gateway.IntentGuildMessages,
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

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s
}

func dummyCommand(e *handler.CommandEvent) error {
	return nil
}
