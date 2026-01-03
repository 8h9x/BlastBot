package friends

import (
	"fmt"

	"github.com/8h9x/BlastBot/internal/sessions"
	"github.com/8h9x/fortgo/usersearch"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var Definition = discord.SlashCommandCreate{
	Name:        "friends",
	Description: "Manage your friends",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionSubCommand{
			Name:        "add",
			Description: "Add a friend.",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:         "id",
					Description:  "The Display Name or User ID of the friend to add.",
					Required:     true,
					Autocomplete: true,
				},
			},
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "list",
			Description: "List all friends.",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "remove",
			Description: "Remove a friend.",
		},
	},
}

func Register(router *handler.Mux) {
	router.Command("/friends/add", addHandler)
	router.Autocomplete("/friends", userLookupAutocompleteHandler)
}

func userLookupAutocompleteHandler(event *handler.AutocompleteEvent) error {
	discordId := event.User().ID

	session, err := sessions.GetSessionForUser(discordId)
	if err != nil {
		return fmt.Errorf("unable to get session for user (%s): %s", discordId, err)
	}

	matches, err := session.UserSearchService.Search(session.CurrentCredentials().AccountID, usersearch.PlatformEpic, event.Data.String("id"))
	if err != nil {
		return fmt.Errorf("unable to search for user (%s): %s", discordId, err)
	}

	var choices = make([]discord.AutocompleteChoice, 0)

	for _, match := range matches[0:25] {
		choices = append(choices, discord.AutocompleteChoiceString{
			Name:  match.Matches[0].Value,
			Value: match.AccountID,
		})
	}

	return event.AutocompleteResult(choices)
}
