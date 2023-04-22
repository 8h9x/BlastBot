package commands

import (
	"blast/api"
	"blast/db"
	"fmt"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/paginator"
)

var vbucks = discord.SlashCommandCreate{
	Name:        "vbucks",
	Description: "Provides information and links to purchase vbucks in your browser.",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:        "amount",
			Description: "The amount of vbucks you want to view.",
			Choices:     VbuckOffers,
			Required:    false,
		},
	},
}

func Vbucks(manager *paginator.Manager) Command {
	return Command{
		Handler: func(event *handler.CommandEvent, blast api.EpicClient, user db.UserEntry, credentials api.UserCredentialsResponse, data discord.SlashCommandInteractionData) error {
			offerID := data.String("amount")

			if offerID != "" {
				offerRes, err := blast.FetchOffers(credentials, offerID)
				if err != nil {
					return err
				}

				offer := offerRes[offerID]

				err = event.CreateMessage(discord.NewMessageCreateBuilder().
					SetEmbeds(discord.NewEmbedBuilder().
						SetColor(0xFB5A32).
						SetTimestamp(time.Now()).
						SetTitle(offer.Title).
						SetDescription(offer.Description).
						SetThumbnail(offer.KeyImages[0].URL).
						AddField("Price (USD)", priceParse(offer.BasePrice), true).
						AddField("Seller", offer.Seller.Name, true).
						AddField("Status", offer.Status, true).
						Build(),
					).
					AddActionRow(discord.NewLinkButton("Purchase", fmt.Sprintf(purchaseLink, offer.ID))).
					Build(),
				)
				if err != nil {
					return err
				}

				return nil
			}

			offerList := make([]string, len(VbuckOffers))

			for _, offer := range VbuckOffers {
				offerList = append(offerList, offer.Value)
			}

			offerRes, err := blast.FetchOffers(credentials, offerList...)
			if err != nil {
				return err
			}

			err = manager.Create(event.Respond, paginator.Pages{
				ID: event.ID().String(),
				PageFunc: func(page int, embed *discord.EmbedBuilder) {
					offer := offerRes[VbuckOffers[page].Value]

					embed.SetTimestamp(time.Now()).
						SetTitle(offer.Title).
						SetDescription(offer.Description).
						SetThumbnail(offer.KeyImages[0].URL).
						AddField("Price (USD)", priceParse(offer.BasePrice), true).
						AddField("Seller", offer.Seller.Name, true).
						AddField("Status", offer.Status, true).
						AddField("Purchase", fmt.Sprintf("**[:link: Epic Store](%s)**", fmt.Sprintf(purchaseLink, offer.ID)), true)
				},
				Pages:      len(VbuckOffers),
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

var purchaseLink = "https://launcher-website-prod07.ol.epicgames.com/purchase?namespace=fn&offers=%s&uePlatform=FNGame"

var VbuckOffers = []discord.ApplicationCommandOptionChoiceString{
	{
		Name:  "100",
		Value: "48a61f0d493942909a529369a66f803b",
	},
	{
		Name:  "200",
		Value: "eca7c7c9c7bc4de69f24d8507cc34ff8",
	},
	{
		Name:  "300",
		Value: "16178da6de794c8d86300a607ec1bda6",
	},
	{
		Name:  "400",
		Value: "a40d8f5f511045aba8be343b87849abd",
	},
	{
		Name:  "500",
		Value: "5f464fcc96764708a54882abc0775e3d",
	},
	{
		Name:  "600",
		Value: "3f62ba5e95f34bf287c06cfd70a8626d",
	},
	{
		Name:  "700",
		Value: "abb65726a6934e8abc2a57e30add0926",
	},
	{
		Name:  "800",
		Value: "5cd3a2123244476499d41ff480ed9340",
	},
	{
		Name:  "900",
		Value: "0d346f08b60b45f1a64cb5c62c8c0a89",
	},
	{
		Name:  "1000",
		Value: "ede05b3c97e9475a8d9be91da65750f0",
	},
	{
		Name:  "1100",
		Value: "ea8bba06cfee427ea0d0b65953438f92",
	},
	{
		Name:  "1200",
		Value: "1c8a96b741484198b6a773b4d11fd048",
	},
	{
		Name:  "1300",
		Value: "c3610a6dc47c40cf980cdbe5722de4a4",
	},
	{
		Name:  "1400",
		Value: "9c586c23738541ba96fc2c458b74558f",
	},
	{
		Name:  "1500",
		Value: "397d63f6f73f4bb1b6628716a08d0c4b",
	},
	{
		Name:  "1600",
		Value: "b207de06940944469b55633f1a8756e2",
	},
	{
		Name:  "1700",
		Value: "233ab6e08b694608a2f82466d94d57d0",
	},
	{
		Name:  "1800",
		Value: "9f85a5d097764566aa0316f86994dd91",
	},
	{
		Name:  "1900",
		Value: "3f1e339fdaa24894a99c1db9f6815db4",
	},
	{
		Name:  "2000",
		Value: "775603a570324885a56762ea81bff788",
	},
	{
		Name:  "2800",
		Value: "559f2ba95f874ec987d0ebfd2cc9c70a",
	},
	{
		Name:  "5000",
		Value: "d900ad5da7ec4eac86918bcfa0c3e698",
	},
	{
		Name:  "13500",
		Value: "4daadb392f1c4ee2b5a3af443e614d2a",
	},
}

func priceParse(num int) string {
	dollar := num / 100
	cents := num % 100
	return fmt.Sprintf("$%d.%02d", dollar, cents)
}
