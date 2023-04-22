package commands

import (
	"blast/api"
	"blast/db"
	"log"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/paginator"
)

var offers = discord.SlashCommandCreate{
	Name:        "offers",
	Description: "Offer catalogue.",
}

func Offers(manager *paginator.Manager) Command {
	return Command{
		Handler: func(event *handler.CommandEvent, blast api.EpicClient, user db.UserEntry, credentials api.UserCredentialsResponse, data discord.SlashCommandInteractionData) error {
			catalogue, err := blast.FetchCatalogue(credentials)
			if err != nil {
				return err
			}

			offers := catalogue.Storefronts[28].CatalogEntries

			// offerList := make([]string, len(offers))
			// for _, offer := range offers {
			// 	offerList = append(offerList, offer.AppStoreID[1])
			// }

			// offerRes, err := blast.FetchOffers(credentials, offerList...)
			// if err != nil {
			// 	return err
			// }

			log.Println(len(catalogue.Storefronts), offers)

			err = manager.Create(event.Respond, paginator.Pages{
				ID: event.ID().String(),
				PageFunc: func(page int, embed *discord.EmbedBuilder) {
					embed.SetDescription("This command is currently under development. Please check back later.")

					// offer := offerRes[offers[page].AppStoreID[1]]

					// embed.SetTimestamp(time.Now()).
					// 	SetTitle(offer.Title).
					// 	SetDescription(offer.Description).
					// 	SetThumbnail(offer.KeyImages[0].URL).
					// 	AddField("Price (USD)", priceParse(offer.BasePrice), true).
					// 	AddField("Seller", offer.Seller.Name, true).
					// 	AddField("Status", offer.Status, true).
					// 	AddField("Purchase", fmt.Sprintf("**[:link: Epic Store](%s)**", fmt.Sprintf(purchaseLink, offer.ID)), true)
				},
				Pages:      len(offers),
				Creator:    event.User().ID,
				ExpireMode: paginator.ExpireModeAfterLastUsage,
			}, false)
			if err != nil {
				return err
			}

			return nil
		},
		LoginRequired:     true,
		EphemeralResponse: false,
	}
}

// schematic editing/upgrading
