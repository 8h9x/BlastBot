package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func (c EpicClient) FetchOffers(credentials UserCredentialsResponse, offerIDs ...string) (OfferResponse, error) {
	values := url.Values{}
	for _, offerID := range offerIDs {
		values.Add("id", offerID)
	}

	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+credentials.AccessToken)

	resp, err := c.Request("GET", fmt.Sprintf("https://catalog-public-service-prod.ol.epicgames.com/catalog/api/shared/bulk/offers?%s", values.Encode()), headers, "")
	if err != nil {
		return OfferResponse{}, err
	}

	defer resp.Body.Close()

	var res OfferResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return OfferResponse{}, err
	}

	return res, nil
}

// https://catalog-public-service-prod.ol.epicgames.com/catalog/api/shared/bulk/offers?id=%7BOfferId%7D
