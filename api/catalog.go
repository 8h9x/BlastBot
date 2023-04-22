package api

import (
	"encoding/json"
	"net/http"
	"time"
)

// https://fortnite-public-service-prod11.ol.epicgames.com/fortnite/api/storefront/v2/catalog

func (c EpicClient) FetchCatalogue(credentials UserCredentialsResponse) (CatalogueResponse, error) {
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Authorization", "Bearer "+credentials.AccessToken)

	resp, err := c.Request("GET", "https://fortnite-public-service-prod11.ol.epicgames.com/fortnite/api/storefront/v2/catalog", headers, "")
	if err != nil {
		return CatalogueResponse{}, err
	}

	defer resp.Body.Close()

	var res CatalogueResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return CatalogueResponse{}, err
	}

	return res, nil
}

type CatalogueResponse struct {
	RefreshIntervalHrs int       `json:"refreshIntervalHrs"`
	DailyPurchaseHrs   int       `json:"dailyPurchaseHrs"`
	Expiration         time.Time `json:"expiration"`
	Storefronts        []struct {
		Name           string `json:"name"`
		CatalogEntries []struct {
			OfferID   string `json:"offerId"`
			DevName   string `json:"devName"`
			OfferType string `json:"offerType"`
			Prices    []struct {
				CurrencyType        string    `json:"currencyType"`
				CurrencySubType     string    `json:"currencySubType"`
				RegularPrice        int       `json:"regularPrice"`
				DynamicRegularPrice int       `json:"dynamicRegularPrice"`
				FinalPrice          int       `json:"finalPrice"`
				SaleExpiration      time.Time `json:"saleExpiration"`
				BasePrice           int       `json:"basePrice"`
			} `json:"prices"`
			Categories   []interface{} `json:"categories"`
			DailyLimit   int           `json:"dailyLimit"`
			WeeklyLimit  int           `json:"weeklyLimit"`
			MonthlyLimit int           `json:"monthlyLimit"`
			Refundable   bool          `json:"refundable"`
			AppStoreID   []string      `json:"appStoreId"`
			Requirements []interface{} `json:"requirements"`
			MetaInfo     []struct {
				Key   string `json:"key,omitempty"`
				Value string `json:"value,omitempty"`
			} `json:"metaInfo"`
			CatalogGroup         string        `json:"catalogGroup"`
			CatalogGroupPriority int           `json:"catalogGroupPriority"`
			SortPriority         int           `json:"sortPriority"`
			Title                string        `json:"title"`
			ShortDescription     string        `json:"shortDescription"`
			Description          string        `json:"description"`
			DisplayAssetPath     string        `json:"displayAssetPath"`
			ItemGrants           []interface{} `json:"itemGrants"`
		} `json:"catalogEntries"`
	} `json:"storefronts"`
}
