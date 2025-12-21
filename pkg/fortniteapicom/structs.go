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

type CosmeticSearchResponse struct {
    Status int `json:"status"`
    Data   CosmeticItem
}