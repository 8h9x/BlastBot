package commands

import (
	"blast/api/consts"
	"blast/db"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"go.mongodb.org/mongo-driver/bson"
)

var mcp = Command{
	Create: discord.SlashCommandCreate{
		Name:        "mcp",
		Description: "Performs a custom MCP request.",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionString{
				Name:        "operation",
				Description: "The name of the operation to perform (ex. QueryProfile)", // unsure yet if case-sensitive.
				Required:    true,
			},
			discord.ApplicationCommandOptionString{
				Name:        "profile",
				Description: "The profile to perform the operation on.",
				Choices: []discord.ApplicationCommandOptionChoiceString{
					{Name: "athena (Battle Royale)", Value: "athena"},
					{Name: "creative", Value: "creative"},
					{Name: "campaign (Save The World)", Value: "campaign"},
					{Name: "common_public (Public Data)", Value: "common_public"},
					{Name: "collections (BR Collections)", Value: "collections"},
					{Name: "common_core (Main Profile)", Value: "common_core"},
					{Name: "metadata (STW Homebase)", Value: "metadata"},
					{Name: "collection_book_people0 (STW Collection Book People)", Value: "collection_book_people0"},
					{Name: "collection_book_schematics0 (STW Collection Book Schematics)", Value: "collection_book_schematics0"},
					{Name: "outpost0 (STW Storage)", Value: "outpost0"},
					{Name: "theater0 (STW Backpack)", Value: "theater0"},
					{Name: "theater1 (STW Event Backpack)", Value: "theater1"},
					{Name: "theater2 (STW Ventures Data)", Value: "theater2"},
					{Name: "recycle_bin (STW Recylce Bin)", Value: "recycle_bin"},
				},
				Required: true,
			},
			discord.ApplicationCommandOptionAttachment{
				Name:        "body",
				Description: "The JSON body to be sent with the request (if ommitted, '{}' will be used).",
			},
		},
	},
	Handler: func(event *events.ApplicationCommandInteractionCreate) error {
		data := event.SlashCommandInteractionData()

		user, err := db.Fetch[db.UserEntry]("users", bson.M{"discordId": event.User().ID.String()})
		if err != nil {
			return err
		}

		refreshCredentials, err := blast.RefreshTokenLogin(consts.FORTNITE_PC_CLIENT_ID, consts.FORTNITE_PC_CLIENT_SECRET, user.Accounts[user.SelectedAccount].RefreshToken)
		if err != nil {
			return err
		}

		operation, profileId, customBodyURL, customBody := data.String("operation"), data.String("profile"), data.Attachment("body").URL, "{}"

		if customBodyURL != "" {
			body, err := fetchJSON(customBodyURL)
			if err != nil {
				return err
			}

			customBody = body
		}

		res, err := blast.ProfileOperationStr(refreshCredentials, operation, profileId, customBody)
		if err != nil {
			return err
		}

		log.Println(res)

		err = event.CreateMessage(discord.NewMessageCreateBuilder().
			AddFile(fmt.Sprintf("%s.%s.blast.json", strings.ToLower(operation), profileId), "Profile operation response.", res).
			Build(),
		)
		if err != nil {
			return err
		}

		res.Close()

		return nil
	},
	LoginRequired: true,
}

// func FetchJSON(url string) (interface{}, error) {
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return "", err
// 	}

// 	defer resp.Body.Close()

// 	var res interface{}
// 	err = json.NewDecoder(resp.Body).Decode(&res)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return res, nil
// }

func fetchJSON(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// https://github.com/LeleDerGrasshalmi/FortniteEndpointsDocumentation/tree/main/Fortnite/MCP

// athena	Battle Royale data
// creative	Creative data
// campaign	Save the World data
// common_public	Public data (e.g. current user banner)
// collections	Battle Royale Collection (Fishing / NPC)
// common_core	Main Profile storing general info like banners, purchases, ban status & history, gift history and much more
// metadata	Save the World Homebase data (e.g info about user permission)
// collection_book_people0	Collectionbook data (Heroes, Survivor, Defender)
// collection_book_schematics0	Collectionbook data (Schematics)
// outpost0	StW Storage
// theater0	StW Backpack
// theater1	Events Backpack (e.g. Frostnite, Hit the Road etc...)
// theater2	Ventures data
// recycle_bin	Stw Recycle Bin
