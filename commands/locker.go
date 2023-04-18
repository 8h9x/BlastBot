package commands

import (
	"blast/api"
	"blast/db"
	"blast/fortnite"
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"log"
	"math"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/fogleman/gg"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var locker = discord.SlashCommandCreate{
	Name:        "locker",
	Description: "Locker stuff.",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionSubCommand{
			Name:        "image",
			Description: "Generates an image of your locker.",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        "type",
					Description: "The type of cosmetic to generate an image for.",
					Choices: []discord.ApplicationCommandOptionChoiceString{
						{Name: "Skins", Value: "outfit"},
						{Name: "Backblings", Value: "backpack"},
						{Name: "Pickaxes", Value: "pickaxe"},
						{Name: "Gliders", Value: "glider"},
						{Name: "Contrails", Value: "contrail"},
						{Name: "Emotes", Value: "emote"},
						{Name: "Emoticons", Value: "emoji"},
						{Name: "Sprays", Value: "spray"},
						{Name: "Toys", Value: "toy"},
						{Name: "Wraps", Value: "wrap"},
						{Name: "Banners", Value: "banner"},
						{Name: "Loading Screens", Value: "loadingscreen"},
						{Name: "Music Packs", Value: "music"},
						// {Name: "Pets", Value: "pet"},
						{Name: "All Cosmetics", Value: "all"},
					},
					Required: true,
				},
				// discord.ApplicationCommandOptionAttachment{
				// 	Name:        "filter",
				// 	Description: "Only includes the cosmetics that match the filter.",
				// },
			},
		},
		// discord.ApplicationCommandOptionSubCommand{
		// 	Name:        "equip",
		// 	Description: "Equips a cosmetic.",
		// 	Options:     opts,
		// },
	},
}

var Locker = Command{
	Handler: func(event *handler.CommandEvent, blast api.EpicClient, user db.UserEntry, credentials api.UserCredentialsResponse, data discord.SlashCommandInteractionData) error {
		startTime := time.Now()
		cosmeticType := data.String("type")
		cosmeticTypeData := itemTypes[cosmeticType]
		// filterTypes := []string{cosmeticTypeData.BackendType}
		// if cosmeticType == "all" {
		// 	filterTypes = []string{"AthenaBackpack", "BannerToken", "AthenaSkyDiveContrail", "AthenaDance", "AthenaGlider", "AthenaLoadingScreen", "AthenaMusicPack", "AthenaCharacter", "AthenaPet", "AthenaPickaxe", "AthenaItemWrap"}
		// }

		athenaProfile, err := blast.ProfileOperationStr(credentials, "QueryProfile", "athena", "{}")
		if err != nil {
			return err
		}

		defer athenaProfile.Close()

		var profile api.AthenaProfile[api.AthenaProfileLockerItem]
		err = json.NewDecoder(athenaProfile).Decode(&profile)
		if err != nil {
			return err
		}

		cosmetics := blast.Fortnite.FetchCosmetics()

		items := profile.ProfileChanges[0].Profile.Items

		filteredCosmetics := make([]fortnite.CosmeticItem, 0)
		for _, value := range items {
			data := cosmetics[strings.Split(value.TemplateID, ":")[1]]

			typeValue := data.Type.Value
			if typeValue == "petcarrier" {
				typeValue = "backpack"
			}

			if cosmeticTypeData.Value == "all" && itemTypes[typeValue].BackendValue != "" {
				filteredCosmetics = append(filteredCosmetics, data)
				continue
			}

			if typeValue == cosmeticTypeData.Value {
				filteredCosmetics = append(filteredCosmetics, data)
			}
		}

		sort.SliceStable(filteredCosmetics, func(i, j int) bool {
			return filteredCosmetics[i].Name < filteredCosmetics[j].Name
		})

		if cosmeticTypeData.Value == "all" {
			sort.SliceStable(filteredCosmetics, func(i, j int) bool {
				prev := filteredCosmetics[i]
				next := filteredCosmetics[j]

				prevSortStr := prev.Type.Value
				if prevSortStr == "petcarrier" {
					prevSortStr = "backpack"
				}

				nextSortStr := next.Type.Value
				if nextSortStr == "petcarrier" {
					nextSortStr = "backpack"
				}

				return typeOrder[prevSortStr] < typeOrder[nextSortStr]
			})
		}

		sort.SliceStable(filteredCosmetics, func(i, j int) bool {
			prev := filteredCosmetics[i]
			next := filteredCosmetics[j]

			prevSortStr := prev.Rarity.BackendValue
			if prev.Series.BackendValue != "" {
				prevSortStr = prev.Series.BackendValue
			}

			nextSortStr := next.Rarity.BackendValue
			if next.Series.BackendValue != "" {
				nextSortStr = next.Series.BackendValue
			}

			return sortOrder[prevSortStr] > sortOrder[nextSortStr]
		})

		count := len(filteredCosmetics)

		// minWidth := 500
		minX := 8.0

		NXx := math.Ceil(math.Sqrt(float64(count)))
		if NXx < minX {
			NXx = minX
		}
		NY := int(math.Ceil(float64(count) / NXx))
		NX := int(NXx)
		width := (64+5)*NX + 5
		// if width < minWidth && count < NY {
		// 	width = minWidth
		// }
		height := (64+10)*NY + 100
		dc := gg.NewContext(width, height)
		dc.SetHexColor("#11111b")
		dc.Clear()

		err = dc.LoadFontFace("assets/fonts/burbank/BurbankBigCondensed-Black.ttf", 36)
		if err != nil {
			return err
		}

		accountInfo, err := blast.FetchAccountInformation(credentials)
		if err != nil {
			return err
		}

		name := accountInfo.DisplayName
		nameW, _ := dc.MeasureString(name)
		x := nameW/2 + 27
		y := 100.0 / 2
		dc.SetRGB(1, 1, 1)
		dc.DrawStringAnchored(name, x, y, 0.5, 0.5)

		slash := fmt.Sprintf(" / %d %ss", count, cases.Title(language.Und).String(cosmeticTypeData.Friendly))
		slashW, _ := dc.MeasureString(slash)
		x = slashW/2 + 27 + nameW
		dc.SetHexColor("#45475a")
		dc.DrawStringAnchored(slash, x, y, 0.5, 0.5)

		logo, err := os.Open("assets/images/blast_100x100.png")
		if err != nil {
			return err
		}

		logoIMG, _, err := image.Decode(logo)
		if err != nil {
			return err
		}

		dc.DrawImage(logoIMG, width-100, 0)

		err = dc.LoadFontFace("assets/fonts/burbank/BurbankBigCondensed-Black.ttf", 8)
		if err != nil {
			return err
		}

		for y := 0; y < NY; y++ {
			for x := 0; x < NX; x++ {
				idx := y*NX + x

				if idx >= len(filteredCosmetics) {
					break
				}

				item := filteredCosmetics[y*NX+x]

				skinIMG, err := os.Open(fmt.Sprintf("assets/images/cosmetics/%s/%s.png", strings.ToLower(item.Type.Value), strings.ToLower(item.ID)))
				if err != nil {
					log.Println(item.ID)
					return err
				}

				defer skinIMG.Close()

				img, _, err := image.Decode(skinIMG)
				if err != nil {
					return err
				}

				w := img.Bounds().Size().X
				h := img.Bounds().Size().Y

				colour := rarityColours[item.Rarity.BackendValue]
				if item.Series.BackendValue != "" {
					colour = rarityColours[item.Series.BackendValue]
				}

				dc.SetRGBA255(colour.R, colour.G, colour.B, 100)
				dc.DrawRoundedRectangle(float64(x*w+(5*x))+5, float64(y*h+100+(10*y)), 64, 64+5, 3)
				dc.Fill()

				dc.DrawImage(img, x*w+(5*x)+5, y*h+100+(10*y))

				dc.SetRGBA255(30, 30, 46, 255) // 128
				dc.DrawRectangle(float64(x*w+(5*x))+5, float64(y*h+100+(10*y))+(64+5-15), 64, 15)
				dc.Fill()

				dc.SetRGBA255(colour.R, colour.G, colour.B, 255)
				dc.DrawRectangle(float64(x*w+(5*x))+5, float64(y*h+100+(10*y))+(64+5-15), 64, 2)
				dc.Fill()

				dc.SetRGBA255(colour.R, colour.G, colour.B, 255)
				dc.SetLineWidth(1.5)
				dc.DrawRoundedRectangle(float64(x*w+(5*x))+5, float64(y*h+100+(10*y)), 64, 64+5, 3)
				dc.Stroke()

				anchorX := float64(x*w+(5*x)) + 64/2
				anchorY := float64(y*h+100+(10*y)) + ((64 + 5) - 15/2)

				dc.SetRGB(1, 1, 1)
				dc.DrawStringAnchored(item.Name, anchorX+5, anchorY, 0.5, 0.5)
			}
		}

		var b bytes.Buffer
		if err := dc.EncodePNG(&b); err != nil {
			return err
		}

		_, err = event.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
			SetEmbeds(discord.NewEmbedBuilder().
				SetTitlef("%s's Locker | %d %ss", accountInfo.DisplayName, count, cases.Title(language.Und).String(cosmeticTypeData.Friendly)).
				SetColor(0xFB5A32).
				SetImage("attachment://image.png").
				SetFooterTextf("Rendered in %v", time.Since(startTime).Round(100*time.Millisecond)).
				SetTimestamp(time.Now()).
				Build()).
			AddFile("image.png", "Profile operation response.", &b).
			Build(),
		)

		return err
	},
	LoginRequired:     true,
	EphemeralResponse: false,
}

type Colour struct {
	R, G, B int
}

type RarityColour struct {
	Background    Colour
	Overlay       Colour
	OverlayBorder Colour
}

var rarityColours = map[string]Colour{
	"EFortRarity::Common":    {R: 205, G: 214, B: 244},
	"EFortRarity::Uncommon":  {R: 166, G: 227, B: 161},
	"EFortRarity::Rare":      {R: 137, G: 180, B: 250},
	"EFortRarity::Epic":      {R: 203, G: 166, B: 247},
	"EFortRarity::Legendary": {R: 250, G: 179, B: 135},
	"ColumbusSeries":         {R: 30, G: 102, B: 245},
	"SlurpSeries":            {R: 129, G: 200, B: 190},
	"ShadowSeries":           {R: 98, G: 104, B: 128},
	"MarvelSeries":           {R: 210, G: 15, B: 57},
	"LavaSeries":             {R: 254, G: 100, B: 11},
	"CreatorCollabSeries":    {R: 82, G: 224, B: 224},
	"PlatformSeries":         {R: 128, G: 120, B: 255},
	"FrozenSeries":           {R: 116, G: 199, B: 236},
	"DCUSeries":              {R: 140, G: 170, B: 238},
	"CUBESeries":             {R: 136, G: 57, B: 239},
}

var sortOrder = map[string]int{
	"EFortRarity::Common":    0,
	"EFortRarity::Uncommon":  1,
	"EFortRarity::Rare":      2,
	"EFortRarity::Epic":      3,
	"EFortRarity::Legendary": 4,
	"ColumbusSeries":         5,
	"SlurpSeries":            6,
	"ShadowSeries":           7,
	"MarvelSeries":           8,
	"LavaSeries":             9,
	"CreatorCollabSeries":    10,
	"PlatformSeries":         11,
	"FrozenSeries":           12,
	"DCUSeries":              13,
	"CUBESeries":             14,
}

type ItemType struct {
	Value        string
	BackendValue string
	Friendly     string
}

var itemTypes = map[string]ItemType{
	"backpack":      {Value: "backpack", BackendValue: "AthenaBackpack", Friendly: "Backbling"},
	"banner":        {Value: "banner", BackendValue: "BannerToken", Friendly: "Banner"},
	"contrail":      {Value: "contrail", BackendValue: "AthenaSkyDiveContrail", Friendly: "Contrail"},
	"emoji":         {Value: "emoji", BackendValue: "AthenaDance", Friendly: "Emoticon"},
	"emote":         {Value: "emote", BackendValue: "AthenaDance", Friendly: "Emote"},
	"glider":        {Value: "glider", BackendValue: "AthenaGlider", Friendly: "Glider"},
	"loadingscreen": {Value: "loadingscreen", BackendValue: "AthenaLoadingScreen", Friendly: "Loading Screen"},
	"music":         {Value: "music", BackendValue: "AthenaMusicPack", Friendly: "Music Pack"},
	"outfit":        {Value: "outfit", BackendValue: "AthenaCharacter", Friendly: "Skin"},
	"pet":           {Value: "petcarrier", BackendValue: "AthenaPetCarrier", Friendly: "Pet"},
	"pickaxe":       {Value: "pickaxe", BackendValue: "AthenaPickaxe", Friendly: "Pickaxe"},
	"spray":         {Value: "spray", BackendValue: "AthenaDance", Friendly: "Spray"},
	"toy":           {Value: "toy", BackendValue: "AthenaDance", Friendly: "Toy"},
	"wrap":          {Value: "wrap", BackendValue: "AthenaVehicleWrap", Friendly: "Wrap"},
	"all":           {Value: "all", BackendValue: "", Friendly: "Cosmetic"},
}

var typeOrder = map[string]int{
	"outfit":        0,
	"backpack":      1,
	"pickaxe":       2,
	"glider":        3,
	"contrail":      4,
	"emote":         5,
	"emoji":         6,
	"spray":         7,
	"toy":           8,
	"wrap":          9,
	"banner":        10,
	"music":         11,
	"loadingscreen": 12,
}

// var backendTypes = map[string]string{
// 	"backpack":      "AthenaBackpack",
// 	"banner":        "BannerToken",
// 	"contrail":      "AthenaSkyDiveContrail",
// 	"emoji":         "AthenaDance", // AthenaDance:Emoji
// 	"emote":         "AthenaDance",
// 	"glider":        "AthenaGlider",
// 	"loadingscreen": "AthenaLoadingScreen",
// 	"music":         "AthenaMusicPack",
// 	"outfit":        "AthenaCharacter",
// 	"pet":           "AthenaPet",
// 	"pickaxe":       "AthenaPickaxe",
// 	"spray":         "AthenaDance", // AthenaDance:SPID
// 	"toy":           "AthenaDance", // AthenaDance:TOY_
// 	"wrap":          "AthenaItemWrap",
// }

func contains(arr []string, str string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}

/*
	{Name: "Backblings", Value: "backpack"},
	{Name: "Banners", Value: "banner"},
	{Name: "Contrails", Value: "contrail"},
	{Name: "Emoticons", Value: "emoji"},
	{Name: "Emotes", Value: "emote"},
	{Name: "Gliders", Value: "glider"},
	{Name: "Loading Screens", Value: "loadingscreen"},
	{Name: "Music Packs", Value: "music"},
	{Name: "Skins", Value: "outfit"},
	{Name: "Pets", Value: "pet"},
	{Name: "Pickaxes", Value: "pickaxe"},
	{Name: "Sprays", Value: "spray"},
	{Name: "Toys", Value: "toy"},
	{Name: "Wraps", Value: "wrap"},
	{Name: "All Cosmetics", Value: "all"},
*/
