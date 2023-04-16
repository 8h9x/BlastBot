package commands

import (
	"blast/api"
	"blast/db"
	"log"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var ephemeral = discord.SlashCommandCreate{
	Name:        "ephemeral",
	Description: "Commands that alter your in-game lobby state for other party members.",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionSubCommand{
			Name:        "crowns",
			Description: "Modifies the amount of crowns you have in the lobby. Can only be seen by other party members.",
		},
	},
}

var EphemeralCrowns = Command{
	Handler: func(event *handler.CommandEvent, blast api.EpicClient, user db.UserEntry, credentials api.UserCredentialsResponse, data discord.SlashCommandInteractionData) error {
		updateData := `{"delete":[],"revision":1,"update":{"Default:AthenaCosmeticLoadout_j":"{\"AthenaCosmeticLoadout\":{\"cosmeticStats\":[{\"statName\":\"TotalVictoryCrowns\",\"statValue\":0},{\"statName\":\"TotalRoyalRoyales\",\"statValue\":\"999\"},{\"statName\":\"HasCrown\",\"statValue\":1}]}}"}`

		party, err := blast.FetchParty(credentials)
		if err != nil {
			return err
		}

		log.Println(party.Current[0].ID)

		res, err := blast.PartyMetaUpdate(credentials, party.Current[0].ID, updateData)
		if err != nil {
			return err
		}

		log.Println(res)

		return nil
	},
	LoginRequired:     true,
	EphemeralResponse: false,
}
