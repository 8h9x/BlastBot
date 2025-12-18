package fortniteapicom

import (
    "net/http"
    "strings"

    "github.com/8h9x/fortgo/request"
)

type Client struct {
    HTTPClient *http.Client
}

func (c *Client) FetchCosmetics() (Cosmetics, error) {
    req, err := request.MakeRequest(
        http.MethodGet,
        "https://fortnite-api.com",
        "v2/cosmetics/br",
    )
    if err != nil {
        return Cosmetics{}, err
    }

    res, err := c.HTTPClient.Do(req)
    if err != nil {
        return Cosmetics{}, err
    }

    resp, err := request.ParseResponse[CosmeticListResponse](res)
    if err != nil {
        return Cosmetics{}, err
    }

    items := make(Cosmetics, len(resp.Data.Data))
    for _, item := range resp.Data.Data {
        items[strings.ToLower(item.ID)] = item
    }

    return items, nil
}