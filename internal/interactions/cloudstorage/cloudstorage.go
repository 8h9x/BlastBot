package cloudstorage

import (
	"fmt"
	"sort"

	"github.com/8h9x/BlastBot/internal/sessions"
	"github.com/8h9x/fortgo"
	"github.com/8h9x/fortgo/fortnite"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var sharedOpts = []discord.ApplicationCommandOption{
	discord.ApplicationCommandOptionString{
		Name:         "file",
		Description:  "File to use.",
		Required:     true,
		Autocomplete: true,
	},
}

var Definition = discord.SlashCommandCreate{
	Name:        "cloudstorage",
	Description: "Collection of commands to manage cloudstorage files.",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionSubCommand{
			Name:        "dump",
			Description: "Export selected .sav file.",
			Options:     sharedOpts,
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "edit",
			Description: "Helper command to modify client settings for selected platform.",
			Options:     sharedOpts,
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "list",
			Description: "List all cloudstorage files.",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "transfer",
			Description: "Copy .sav file from one platform to another",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "upload",
			Description: "Replace cloudstorage file on the server with uploaded one.",
		},
	},
}

func Register(router *handler.Mux) {
	router.Command("/cloudstorage/list", listHandler)
	router.Autocomplete("/cloudstorage", autocompleteHandlerTEST)
}

func listUserCloudstorageSorted(session *fortgo.Client) ([]fortnite.CloudstorageFilePointerUser, error) {
	accountID := session.CurrentCredentials().AccountID

	cloudstorageList, err := session.FortniteService.ListUserCloudstorage(accountID)
	if err != nil {
		return nil, fmt.Errorf("unable to get cloudstorage files for account (%s): %s", accountID, err)
	}

	sort.Slice(cloudstorageList, func(i, j int) bool {
		iFilename := cloudstorageList[i].Filename
		jFilename := cloudstorageList[j].Filename

		if isUUID(iFilename) && !isUUID(jFilename) {
			return false
		}
		if !isUUID(iFilename) && isUUID(jFilename) {
			return true
		}

		return cloudstorageList[i].Filename < cloudstorageList[j].Filename
	})

	return cloudstorageList, nil
}

func autocompleteHandlerTEST(event *handler.AutocompleteEvent) error {
	discordId := event.User().ID

	session, err := sessions.GetSessionForUser(discordId)
	if err != nil {
		return fmt.Errorf("unable to get session for user (%s): %s", discordId, err)
	}

	var choices = make([]discord.AutocompleteChoice, 0)

	cloudstorageList, err := listUserCloudstorageSorted(session)
	if err != nil {
		return err
	}

	for _, file := range cloudstorageList {
		choices = append(choices, discord.AutocompleteChoiceString{
			Name:  file.Filename,
			Value: file.Filename,
		})
	}

	return event.AutocompleteResult(choices)
}
