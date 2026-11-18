package fortniteapicom

import "time"

//type Cosmetics map[string]CosmeticItem

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
		SmallIcon string `json:"smallIcon"`
		Icon      string `json:"icon"`
		Featured  any    `json:"featured"`
		Other     any    `json:"other"`
	} `json:"images"`
	Variants            any       `json:"variants"`
	SearchTags          any       `json:"searchTags"`
	GameplayTags        []string  `json:"gameplayTags"`
	MetaTags            any       `json:"metaTags"`
	ShowcaseVideo       any       `json:"showcaseVideo"`
	DynamicPakID        any       `json:"dynamicPakId"`
	ItemPreviewHeroPath string    `json:"itemPreviewHeroPath"`
	DisplayAssetPath    any       `json:"displayAssetPath"`
	DefinitionPath      any       `json:"definitionPath"`
	Path                string    `json:"path"`
	Added               time.Time `json:"added"`
	ShopHistory         any       `json:"shopHistory"`
}

type CosmeticSearchResponse struct {
	Status int `json:"status"`
	Data   CosmeticItem
}
