package fortnite

import (
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/nfnt/resize"
)

type Fortnite struct {
	cosmetics Cosmetics
}

func New() (*Fortnite, error) {
	cosmetics, err := fetchCosmetics()

	return &Fortnite{
		cosmetics: cosmetics,
	}, err
}

func (f *Fortnite) FetchCosmetics() Cosmetics {
	return f.cosmetics
}

func (f *Fortnite) FetchCosmetic(ID string) CosmeticItem {
	cosmetics := f.FetchCosmetics()

	return cosmetics[ID]
}

func (f *Fortnite) FetchCosmeticByName(name string) CosmeticItem {
	cosmetics := f.FetchCosmetics()

	for _, item := range cosmetics {
		if item.Name == name {
			return item
		}
	}
	return CosmeticItem{}
}

func (f *Fortnite) DownloadCosmeticIcons() error {
	cosmetics := f.FetchCosmetics()

	ensureExists("assets")
	ensureExists("assets", "images")
	ensureExists("assets", "images", "cosmetics")

	for _, item := range cosmetics {
		filePath := fmt.Sprintf("assets/images/cosmetics/%s/%s.png", strings.ToLower(item.Type.Value), strings.ToLower(item.ID))

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			log.Println("Downloading", item.Name, "to", filePath)

			ensureExists("assets", "images", "cosmetics", item.Type.Value)

			err = downloadImage(item, filePath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func fetchCosmetics() (Cosmetics, error) {
	resp, err := http.Get("https://fortnite-api.com/v2/cosmetics/br")
	if err != nil {
		return Cosmetics{}, err
	}

	defer resp.Body.Close()

	var res CosmeticListResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return Cosmetics{}, err
	}

	items := make(Cosmetics)
	for _, item := range res.Data {
		items[strings.ToLower(item.ID)] = item
	}

	return items, nil
}

type Cosmetics map[string]CosmeticItem

type CosmeticListResponse struct {
	Status int `json:"status"`
	Data   []CosmeticItem
}

type CosmeticItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        struct {
		Value        string `json:"value"`
		DisplayValue string `json:"displayValue"`
		BackendValue string `json:"backendValue"`
	} `json:"type"`
	Rarity struct {
		Value        string `json:"value"`
		DisplayValue string `json:"displayValue"`
		BackendValue string `json:"backendValue"`
	} `json:"rarity"`
	Series struct {
		BackendValue string
		Colors       []string
		Image        string
		Value        string
	} `json:"series"`
	Set struct {
		Value        string `json:"value"`
		Text         string `json:"text"`
		BackendValue string `json:"backendValue"`
	} `json:"set"`
	Introduction struct {
		Chapter      string `json:"chapter"`
		Season       string `json:"season"`
		Text         string `json:"text"`
		BackendValue int    `json:"backendValue"`
	} `json:"introduction"`
	Images struct {
		SmallIcon string      `json:"smallIcon"`
		Icon      string      `json:"icon"`
		Featured  interface{} `json:"featured"`
		Other     interface{} `json:"other"`
	} `json:"images"`
	Variants            interface{} `json:"variants"`
	SearchTags          interface{} `json:"searchTags"`
	GameplayTags        []string    `json:"gameplayTags"`
	MetaTags            interface{} `json:"metaTags"`
	ShowcaseVideo       interface{} `json:"showcaseVideo"`
	DynamicPakID        interface{} `json:"dynamicPakId"`
	ItemPreviewHeroPath string      `json:"itemPreviewHeroPath"`
	DisplayAssetPath    interface{} `json:"displayAssetPath"`
	DefinitionPath      interface{} `json:"definitionPath"`
	Path                string      `json:"path"`
	Added               time.Time   `json:"added"`
	ShopHistory         interface{} `json:"shopHistory"`
}

func downloadImage(item CosmeticItem, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	imageURL := item.Images.Icon
	if imageURL == "" {
		imageURL = item.Images.SmallIcon
	}

	log.Println(item.Images, imageURL)

	response, err := http.Get(imageURL)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return err
	}

	defer response.Body.Close()

	img, _, err := image.Decode(response.Body)
	if err != nil {
		return err
	}

	const size = 64

	resized := resize.Resize(size, size, img, resize.MitchellNetravali)

	err = png.Encode(file, resized)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func ensureExists(dirs ...string) {
	dir := ""
	for _, d := range dirs {
		dir = path.Join(dir, d)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0644)
	}
}
